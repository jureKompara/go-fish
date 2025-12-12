package main

import (
	"go-fish/internal/engine"
	"go-fish/internal/eval"
)

func Q(p *engine.Position, alpha, beta int) int {
	qNodes++
	us := p.ToMove

	// ---------------------------
	// Case 1: side to move is in check → full evasion search
	// ---------------------------
	if p.InCheck(us) {
		best := -INF
		moves := p.Movebuff[p.Ply][:0]
		p.GenMoves(&moves)

		foundLegal := false

		for _, m := range moves {
			p.Make(m)
			if p.InCheck(us) {
				p.Unmake(m)
				continue
			}
			foundLegal = true

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

		if !foundLegal { // checkmate
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

	moves := p.Movebuff[p.Ply][:0]
	//TODO: implement GenTactics!!!!!!!!!!!!!!!!
	//p.GenTactics(&moves)

	for _, m := range moves {
		p.Make(m)
		if p.InCheck(us) {
			p.Unmake(m)
			continue
		}

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
