package main

import (
	"go-fish/internal/engine"
	"go-fish/internal/eval"
	"sort"
)

func Q(p *engine.Position, alpha, beta int32) int32 {

	if p.Is3Fold() {
		return 0
	}
	qNodes++

	checkers := p.Checkers(p.Kings[p.Stm], p.Stm^1)

	if checkers != 0 {
		// ---------------------------
		// Case 1: side to move is in check → full evasion search
		// ---------------------------
		best := -INF
		moves := p.GenEvasions(checkers)

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

		if len(moves) == 0 { // checkmate
			return -MATE + int32(p.Ply)
		}
		return best
	}
	// ---------------------------
	// Case 2: not in check → normal quiescence
	// ---------------------------

	// stand pat: "if we do nothing"
	stand := eval.Pst(p)

	if stand >= beta {
		return stand
	}
	if stand > alpha {
		alpha = stand
	}

	moves := p.GenTactics()

	sort.Slice(moves, func(i, j int) bool {
		return engine.MvvLvaScore(p, moves[i]) > engine.MvvLvaScore(p, moves[j])
	})

	for _, m := range moves {

		//what we gain from the capture
		gain := eval.Points[p.Board[m.To()]]

		// delta prune
		if stand+gain+50 < alpha {
			continue
		}

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
