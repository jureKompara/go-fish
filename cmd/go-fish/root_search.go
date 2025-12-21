package main

import (
	"go-fish/internal/engine"
)

const INF int32 = 1000000000
const MATE int32 = 1000000

func RootSearch(p *engine.Position, depth int) engine.Move {

	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]

	if n == 0 {
		return 0
	}

	prev := int32(0)
	const base int32 = 25

	for d := 1; d <= depth; d++ {
		w := base
		a := prev - w
		b := prev + w

		for {
			bestScore := -INF
			bestIdx := 0

			for i, m := range moves {
				p.Make(m)
				score := -AB(p, -b, -a, d-1) // root aspiration window
				p.Unmake(m)

				if score > bestScore {
					bestScore = score
					bestIdx = i
				}
			}

			// fail-low
			if bestScore <= a {
				a -= w
				w *= 2
				continue
			}
			// fail-high
			if bestScore >= b {
				b += w
				w *= 2
				continue
			}

			prev = bestScore
			moves[0], moves[bestIdx] = moves[bestIdx], moves[0]
			break
		}
	}
	return moves[0]

}
