package eval

import (
	"go-fish/internal/engine"
	"math/bits"
)

var points = [6]int{100, 300, 300, 500, 900, 0}

type EvalFunc func(*engine.Position) int

func Material(p *engine.Position) int {
	score := 0
	for color := range 2 {
		for piece := engine.PAWN; piece <= engine.KING; piece++ {
			score += bits.OnesCount64(p.PieceBB[color][piece]) * points[piece] * (-2*color + 1)
		}
	}
	return score * (-2*p.ToMove + 1)
}
