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

		for _, m := range moves {
			score := int32(0)
			p.Make(m)
			if !p.Is3Fold() {
				// null-window
				score = -Q(p, -alpha-1, -alpha)
				if score > alpha && score < beta {
					// re-search
					score = -Q(p, -beta, -alpha)
				}
			}
			p.Unmake(m)

			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				return beta
			}
		}
		return alpha
	}

	// ---------------------------
	// Case 2: not in check → normal quiescence
	// ---------------------------

	// stand pat: "if we do nothing"
	stand := eval.Eval(p)

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

	//put best capture first
	if capCount > 1 {
		best := 0
		bestScore := engine.MvvLvaScore(p, moves[0])
		for i := 1; i < capCount; i++ {
			s := engine.MvvLvaScore(p, moves[i])
			if s > bestScore {
				bestScore = s
				best = i
			}
		}
		moves[0], moves[best] = moves[best], moves[0]
	}

	//capture loop
	for _, m := range moves[:capCount] {
		//what we gain from the capture
		gain := int32(0)

		if m.IsEP() {
			gain = eval.PawnValue
		} else {
			gain = eval.Points[p.Board[m.To()]]
			if m.IsPromo() {
				gain += eval.Points[m.Promo()] - eval.PawnValue
			}
		}

		// delta prune
		if stand+gain+50 < alpha {
			continue
		}

		p.Make(m)
		var score int32
		// null-window
		score = -Q(p, -alpha-1, -alpha)
		if score > alpha && score < beta {
			// re-search
			score = -Q(p, -beta, -alpha)
		}

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
		gain := eval.Points[m.Promo()] - eval.PawnValue

		// delta prune
		if stand+gain+50 < alpha {
			continue
		}
		p.Make(m)
		var score int32
		// null-window
		score = -Q(p, -alpha-1, -alpha)
		if score > alpha && score < beta {
			// re-search
			score = -Q(p, -beta, -alpha)
		}
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
