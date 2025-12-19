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

	moves := p.Movebuff[p.Ply][:0]
	scores := [256]int{}
	n := p.GenMoves(moves)
	moves = moves[:n]

	for i, m := range moves {
		p.Make(m)
		score := -AlphaBeta(p, -INF, INF, 0)
		scores[i] = score
		p.Unmake(m)
		if score > bestScore {
			bestScore = score
			bestMove = m
		}
	}

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
