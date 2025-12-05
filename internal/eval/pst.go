package eval

import (
	"go-fish/internal/engine"
	"math/bits"
)

func Pst(p *engine.Position) int {
	score := 0
	for color := range 2 {
		for piece := engine.PAWN; piece <= engine.KING; piece++ {
			score += bits.OnesCount64(p.PieceBB[color*6+piece]) * points[piece] * (-2*color + 1)
		}
	}
	return score
}
