package main

import (
	"go-fish/internal/engine"
	"go-fish/internal/eval"
)

func Q(p *engine.Position, alpha, beta int) int {
	qNodes++
	// ---------------------------
	// Case 1: side to move is in check → full evasion search
	// ---------------------------
	if p.InCheck() {
		best := -INF
		moves := p.Movebuff[p.Ply][:]
		n := p.GenMoves(moves)
		moves = moves[:n]

		for _, m := range moves {
			p.Make(m)
			score := -Q(p, -beta, -alpha)
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

		if n == 0 { // checkmate
			return -MATE + p.Ply
		}
		return best
	}

	if p.Ply >= 8 {
		return eval.Pst(p)
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
	//TODO: implement GenTactics!!!!!!!!!!!!!!!!
	n := p.GenMoves(moves)
	moves = moves[:n]

	for _, m := range moves {
		p.Make(m)
		score := -Q(p, -beta, -alpha)
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
