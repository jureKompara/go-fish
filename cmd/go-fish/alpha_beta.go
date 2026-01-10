package main

import (
	"go-fish/internal/engine"
)

var abNodes uint64
var qNodes uint64

const (
	EXACT uint8 = iota
	LOWER
	UPPER
)

func AB(p *engine.Position, alpha, beta int32, depth int) int32 {
	abNodes++

	//repetition check
	if p.Is3Fold() {
		return 0
	}

	//something something checkmate
	alpha = max(alpha, -MATE+int32(p.Ply))
	beta = min(beta, MATE-int32(p.Ply))

	if alpha >= beta {
		return alpha
	}

	//we start quiesence at leaf nodes
	if depth == 0 {
		return Q(p, alpha, beta)
	}

	//TT probe
	bucket := &engine.TT[p.Hash&engine.IndexMask]

	entry := bucket.Probe(p.Hash)

	TTProbe++

	if entry != nil && entry.Depth >= uint8(depth) {
		ttScore := loadScore(entry.Score, p.Ply)
		TTHit++
		switch entry.BoundType {
		case EXACT:
			return ttScore

		case LOWER:
			alpha = max(alpha, ttScore)

		case UPPER:
			beta = min(beta, ttScore)
		}

		if alpha >= beta {
			ttCutoffs++
			return ttScore
		}
	}

	// Null-move pruning
	if depth >= 3 && !p.InCheck() &&
		(p.PieceBB[p.Stm][engine.KNIGHT]|
			p.PieceBB[p.Stm][engine.BISHOP]|
			p.PieceBB[p.Stm][engine.ROOK]|
			p.PieceBB[p.Stm][engine.QUEEN]) != 0 {

		R := min(2+depth/4, depth-1)

		p.MakeNull()

		// Null search uses a null window at beta
		score := -AB(p, -beta, -beta+1, depth-1-R)

		p.UnmakeNull()

		if score >= beta {
			return beta // fail-high cutoff
		}
	}

	moves := p.GenMoves()
	n := len(moves)

	originalAlpha := alpha
	originalBeta := beta
	best := -INF
	var bestMove engine.Move

	if n > 0 {
		var score int32

		//hash move is first!
		off := 0
		if entry != nil {
			hashMove := entry.HashMove
			for i := range moves {
				if moves[i] == hashMove {
					off = 1
					moves[0], moves[i] = moves[i], moves[0]
					m := moves[0]

					p.Make(m)
					score = -AB(p, -beta, -alpha, depth-1)
					p.Unmake(m)

					if score > best {
						best = score
						bestMove = m
					}
					if score > alpha {
						alpha = score
					}
					if alpha >= beta {
						if !m.IsCapture() {
							engine.History[p.Stm][m.From()][m.To()] += depth * depth
							k0 := engine.Killers[p.Ply][0]
							if k0 != m {
								engine.Killers[p.Ply][1] = k0
								engine.Killers[p.Ply][0] = m
							}
						}
						goto Jmp
					}
					break
				}
			}
		}

		//partition captures first
		write := off + headCaptures(moves[off:])

		//capture loop
		for i := off; i < write; i++ {
			bestIdx := i
			bestScore := engine.MvvLvaScore(p, moves[i])
			for j := i + 1; j < write; j++ {
				s := engine.MvvLvaScore(p, moves[j])
				if s > bestScore {
					bestScore = s
					bestIdx = j
				}
			}
			moves[i], moves[bestIdx] = moves[bestIdx], moves[i]
			m := moves[i]

			p.Make(m)
			if i == 0 {
				//full-window
				score = -AB(p, -beta, -alpha, depth-1)
			} else {
				// null-window
				score = -AB(p, -alpha-1, -alpha, depth-1)
				if score > alpha && score < beta {
					// re-search
					score = -AB(p, -beta, -alpha, depth-1)
				}
			}
			p.Unmake(m)

			if score > best {
				best = score
				bestMove = m
			}
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				goto Jmp
			}
		}

		quietsStart := write
		//put killers first
		k := engine.Killers[p.Ply][0]
		for i := write; i < n; i++ {
			if k == moves[i] {
				moves[write], moves[i] = moves[i], moves[write]
				write++
				break
			}
		}

		k = engine.Killers[p.Ply][1]
		for i := write; i < n; i++ {
			if k == moves[i] {
				moves[write], moves[i] = moves[i], moves[write]
				write++
				break
			}
		}

		//put best history first...ignore the rest
		if n-write > 1 {
			to0 := moves[write].To()
			fr0 := moves[write].From()
			bst := engine.History[p.Stm][fr0][to0]
			bstI := write
			for i := write; i < n; i++ {
				toi := moves[i].To()
				fri := moves[i].From()
				h := engine.History[p.Stm][fri][toi]
				if h > bst {
					bst = h
					bstI = i
				}
			}
			moves[write], moves[bstI] = moves[bstI], moves[write]
		}

		//quiets loop
		for i := quietsStart; i < n; i++ {
			m := moves[i]
			p.Make(m)
			if i == 0 {
				score = -AB(p, -beta, -alpha, depth-1)
			} else {
				// null-window
				score = -AB(p, -alpha-1, -alpha, depth-1)
				if score > alpha && score < beta {
					// re-search
					score = -AB(p, -beta, -alpha, depth-1)
				}
			}
			p.Unmake(m)

			if score > best {
				best = score
				bestMove = m
			}
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				engine.History[p.Stm][m.From()][m.To()] += depth * depth
				k0 := engine.Killers[p.Ply][0]
				if k0 != m {
					engine.Killers[p.Ply][1] = k0
					engine.Killers[p.Ply][0] = m
				}
				break
			}
		}

	} else if p.InCheck() {
		//checkmate
		best = -MATE + int32(p.Ply)
	} else {
		//stalemate
		best = 0
	}

Jmp:

	boundType := EXACT
	if best <= originalAlpha {
		boundType = UPPER
	} else if best >= originalBeta {
		boundType = LOWER
	}

	bucket.Store(entry, &engine.TTEntry{
		Hash:      p.Hash,
		Score:     storeScore(best, p.Ply),
		HashMove:  bestMove,
		Depth:     uint8(depth),
		BoundType: boundType,
	})

	return best
}

func storeScore(score int32, ply int) int32 {
	if score > 10000 {
		return score + int32(ply)
	}
	if score < -10000 {
		return score - int32(ply)
	}
	return score
}

func loadScore(score int32, ply int) int32 {
	if score > 10000 {
		return score - int32(ply)
	}
	if score < -10000 {
		return score + int32(ply)
	}
	return score
}
