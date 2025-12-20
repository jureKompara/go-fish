package main

import (
	"go-fish/internal/engine"
)

const INF = 1000000000
const MATE = 1000000

func RootSearch(p *engine.Position, depth int) engine.Move {
	bestScore := -INF

	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]

	if n == 0 {
		return 0
	}

	bestIdx := 0
	for d := range depth {
		bestScore = -INF
		bestIdx = 0
		for i, m := range moves {
			p.Make(m)
			score := -AB(p, -INF, INF, d)
			p.Unmake(m)
			if score > bestScore {
				bestScore = score
				bestIdx = i
			}
		}
		moves[0], moves[bestIdx] = moves[bestIdx], moves[0]
	}
	return moves[0]
}
