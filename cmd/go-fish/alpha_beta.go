package main

import (
	"go-fish/internal/engine"
	"sort"
)

var abNodes uint64
var qNodes uint64

const (
	EXACT uint8 = iota
	LOWER
	UPPER
)

func AB(p *engine.Position, alpha, beta int32, depth int) int32 {

	if is3Fold(p) {
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

	hit := entry.Hash == p.Hash
	if hit && entry.Depth >= uint8(depth) {
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
		return Q(p, alpha, beta, QDEPTH)
	}

	originalAlpha := alpha
	originalBeta := beta

	abNodes++

	best := -INF
	var bestMove engine.Move

	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]

	//hash move is first!
	off := 0
	if p.Hash == entry.Hash {
		hashMove := entry.HashMove
		for i, m := range moves {
			if m == hashMove {
				off = 1
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	//partition captures first
	write := off
	for i := off; i < n; i++ {
		if engine.IsCapture(moves[i].Flags()) {
			moves[write], moves[i] = moves[i], moves[write]
			write++
		}
	}

	//only sort if there is more than 1 capture
	if write-off > 1 {
		sort.Slice(moves[off:write], func(i, j int) bool {
			return engine.MvvLvaScore(p, moves[off+i]) > engine.MvvLvaScore(p, moves[off+j])
		})
	}

	first := true
	for _, m := range moves {
		p.Make(m)
		var score int32
		if first {
			score = -AB(p, -beta, -alpha, depth-1)
			first = false
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
			break
		}
	}

	if n == 0 {
		if p.InCheck() {
			best = -MATE + int32(p.Ply)
		} else {
			best = 0
		}
	}

	boundType := EXACT
	if best <= originalAlpha {
		boundType = UPPER
	} else if best >= originalBeta {
		boundType = LOWER
	}

	if entry.Hash == p.Hash {
		// same position: always refresh best move (and age if you have it)
		entry.HashMove = bestMove

		// only overwrite score/depth/flag if new info is at least as deep
		if uint8(depth) < entry.Depth {
			engine.TT[index].Score = best
			engine.TT[index].BoundType = boundType
			return best
		}
		// else fallthrough to full overwrite
	} else if entry.Hash != 0 && entry.Depth > uint8(depth) {
		// collision and old is deeper: keep old
		return best
	}

	engine.TT[index] = engine.TTEntry{
		Hash:      p.Hash,
		Depth:     uint8(depth),
		Score:     storeScore(best, p.Ply),
		BoundType: boundType,
		HashMove:  bestMove}

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

func is3Fold(p *engine.Position) bool {
	if p.HalfMove >= 8 {
		count := 0
		for i := p.Ply - 2; i >= max(0, p.Ply-p.HalfMove); i -= 2 {
			if p.Hash == p.HashHistory[i] {
				count++
				if count == 2 {
					return true
				}
			}
		}
	}
	return false
}
