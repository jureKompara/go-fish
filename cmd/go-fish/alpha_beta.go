package main

import (
	"go-fish/internal/engine"
	"go-fish/internal/eval"
)

var abNodes uint64
var qNodes uint64

const (
	EXACT uint8 = iota
	LOWER
	UPPER
)

func AB(p *engine.Position, alpha, beta, depth int) int {

	//TT probe
	TTProbe++
	entry := engine.TT[p.Hash&engine.IndexMask]

	if entry.Hash == p.Hash && entry.Depth >= uint8(depth) {
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

	originalAlpha := alpha
	originalBeta := beta

	//we start quiesence at leaf nodes
	if depth == 0 {
		return eval.Pst(p)
	}
	abNodes++

	best := -INF
	var bestMove engine.Move

	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]

	if p.Hash == entry.Hash {
		hashMove := entry.BestMove
		for i, m := range moves {
			if m == hashMove {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	for _, m := range moves {
		p.Make(m)
		score := -AB(p, -beta, -alpha, depth-1)
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
			best = -MATE + p.Ply
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

	engine.TT[p.Hash&engine.IndexMask] = engine.TTEntry{
		Hash:      p.Hash,
		Depth:     uint8(depth),
		Score:     storeScore(best, p.Ply),
		BoundType: boundType,
		BestMove:  bestMove}

	return best
}

func storeScore(score, ply int) int {
	if score > 10000 { // winning mate for side to move
		return score + ply
	}
	if score < -10000 { // losing mate
		return score - ply
	}
	return score
}

func loadScore(score, ply int) int {
	if score > 10000 {
		return score - ply
	}
	if score < -10000 {
		return score + ply
	}
	return score
}
