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

	alpha = max(alpha, -MATE+int32(p.Ply))
	beta = min(beta, MATE-int32(p.Ply))

	if alpha >= beta {
		return alpha
	}

	//TT probe
	index := p.Hash & engine.IndexMask
	entry := engine.TT[index]
	TTProbe++

	isHit := entry.Hash == p.Hash
	if isHit && entry.Depth >= uint8(depth) {
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

	//we start quiesence at leaf nodes
	if depth == 0 {
		return Q(p, alpha, beta)
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
		if p.Hash == entry.Hash {
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
						k0 := engine.Killers[p.Ply][0]
						if !engine.IsCapture(m.Flags()) && k0 != m {
							engine.Killers[p.Ply][1] = k0
							engine.Killers[p.Ply][0] = m
							engine.History[p.Stm][m.From()][m.To()] += depth * depth
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

		//sorts quiets with history
		for i := write; i < n; i++ {
			ito := moves[i].To()
			ifr := moves[i].From()
			bst := engine.History[p.Stm][ifr][ito]
			bstI := i
			for j := i + 1; j < n; j++ {
				jto := moves[j].To()
				jfr := moves[j].From()
				hj := engine.History[p.Stm][jfr][jto]
				if hj > bst {
					bst = hj
					bstI = j
				}
			}
			moves[i], moves[bstI] = moves[bstI], moves[i]
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
				k0 := engine.Killers[p.Ply][0]
				if k0 != m {
					engine.Killers[p.Ply][1] = k0
					engine.Killers[p.Ply][0] = m
					engine.History[p.Stm][m.From()][m.To()] += depth * depth
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
	if isHit && entry.Depth > uint8(depth) {
		return best
	}

	boundType := EXACT
	if best <= originalAlpha {
		boundType = UPPER
	} else if best >= originalBeta {
		boundType = LOWER
	}

	engine.TT[index] = engine.TTEntry{
		Hash:      p.Hash,
		Depth:     uint8(depth),
		Score:     storeScore(best, p.Ply),
		BoundType: boundType,
		HashMove:  bestMove,
	}

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
