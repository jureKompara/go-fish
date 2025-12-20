package main

import (
	"go-fish/internal/engine"
	"go-fish/internal/eval"
)

const QDEPTH = 4

func Q(p *engine.Position, alpha, beta, qDepth int) int {

	//max qDepth-> static eval
	if qDepth == 0 {
		return eval.Pst(p)
	}
	qNodes++

	checkers := p.Checkers(p.Kings[p.Stm], p.Stm^1)

	if checkers != 0 {
		// ---------------------------
		// Case 1: side to move is in check → full evasion search
		// ---------------------------
		best := -INF
		moves := p.Movebuff[p.Ply][:]
		n := p.GenEvasions(moves, checkers)
		moves = moves[:n]

		for _, m := range moves {
			p.Make(m)
			score := -Q(p, -beta, -alpha, qDepth-1)
			p.Unmake(m)

			if score > best {
				best = score
			}
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				return alpha
			}
		}

		if n <= 0 { // checkmate
			return -MATE + p.Ply
		}
		return best
	}
	// ---------------------------
	// Case 2: not in check → normal quiescence
	// ---------------------------

	// stand pat: "if we do nothing"
	staticEval := eval.Pst(p)

	if staticEval >= beta {
		return staticEval
	}
	if staticEval > alpha {
		alpha = staticEval
	}

	moves := p.Movebuff[p.Ply][:]
	n := p.GenTactics(moves)
	moves = moves[:n]

	for _, m := range moves {
		p.Make(m)
		score := -Q(p, -beta, -alpha, qDepth-1)
		p.Unmake(m)

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
