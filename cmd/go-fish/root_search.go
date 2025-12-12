package main

import (
	"go-fish/internal/engine"
	"sort"
)

const INF = 1000000000
const MATE = 1000000

func RootSearch(p *engine.Position, depth int) engine.Move {
	bestScore := -INF
	var bestMove engine.Move
	us := p.ToMove

	moves := p.Movebuff[p.Ply][:0]
	legalCount := 0
	scores := [256]int{}
	p.GenMoves(&moves)

	for _, m := range moves {
		p.Make(m)
		if p.InCheck(us) {
			p.Unmake(m)
			continue
		}
		moves[legalCount] = m

		score := -AlphaBeta(p, -INF, INF, 0)
		scores[legalCount] = score
		legalCount++
		p.Unmake(m)
		if score > bestScore {
			bestScore = score
			bestMove = m
		}
	}

	moves = p.Movebuff[p.Ply][:legalCount]

	for d := 1; d < depth; d++ {
		bestScore = -INF
		sort.Slice(moves, func(i, j int) bool {
			return scores[i] > scores[j]
		})

		for i, m := range moves {
			p.Make(m)

			score := -AlphaBeta(p, -INF, INF, d)
			scores[i] = score
			p.Unmake(m)
			if score > bestScore {
				bestScore = score
				bestMove = m
			}
		}
	}
	return bestMove
}
