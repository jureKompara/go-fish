package main

import (
	"go-fish/internal/engine"
	"go-fish/internal/eval"
)

func Q(p *engine.Position, alpha, beta int32) int32 {
	qNodes++

	checkers := p.Checkers()

	if checkers != 0 {
		// ---------------------------
		// Case 1: side to move is in check → full evasion search
		// ---------------------------
		moves := p.GenEvasions(checkers)

		if len(moves) == 0 { // checkmate
			return -MATE + int32(p.Ply)
		}

		headCaptures(moves)

		best := -INF

		for _, m := range moves {
			score := int32(0)
			p.Make(m)
			if !p.Is3Fold() {
				score = -Q(p, -beta, -alpha)
			}
			p.Unmake(m)

			if score > best {
				best = score
			}
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				return beta
			}
		}
		return best
	}
	// ---------------------------
	// Case 2: not in check → normal quiescence
	// ---------------------------

	// stand pat: "if we do nothing"
	stand := eval.Pst(p)

	if stand >= beta {
		return beta
	}
	if stand > alpha {
		alpha = stand
	}

	//captures + promotions
	moves := p.GenTactics()

	//pushed promos to the end
	//and returns captures count
	capCount := tailQuiets(moves)

	for i := range capCount {
		best := i
		bestScore := engine.MvvLvaScore(p, moves[i])
		for j := i + 1; j < capCount; j++ {
			s := engine.MvvLvaScore(p, moves[j])
			if s > bestScore {
				bestScore = s
				best = j
			}
		}
		moves[i], moves[best] = moves[best], moves[i]

		m := moves[i]

		//what we gain from the capture
		gain := int32(0)

		flags := m.Flags()
		if flags == engine.EP {
			gain = eval.PawnValue
		} else {
			gain = eval.Points[p.Board[m.To()]]
			if engine.IsPromo(flags) {
				gain += eval.Points[engine.Promo(flags)] - eval.PawnValue
			}
		}

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

	//try the promotions
	for _, m := range moves[capCount:] {
		//what we gain from the capture
		gain := eval.Points[engine.Promo(m.Flags())] - eval.PawnValue

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
