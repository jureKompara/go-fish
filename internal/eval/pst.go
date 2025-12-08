package eval

import (
	"go-fish/internal/engine"
)

func Pst(p *engine.Position) int {
	score := 0
	for color := range 2 {
		sign := -2*color + 1
		for piece := engine.PAWN; piece <= engine.KING; piece++ {
			material := points[piece]
			bb := p.PieceBB[piece+color*6]
			for bb != 0 {
				sq := engine.PopLSB(&bb)
				score += sign * (material + pst[piece][sq^56*(1-color)])
			}
		}
	}
	return score * (-2*p.ToMove + 1)
}
