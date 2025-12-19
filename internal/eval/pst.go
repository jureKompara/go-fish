package eval

import (
	"go-fish/internal/engine"
)

var points = [6]int{100, 310, 320, 500, 900, 0}

// returns the evaluation of the position from side to move POV
// counts material and PST to get the eval
func Pst(p *engine.Position) int {
	score := 0
	for color := range 2 {
		sign := -2*color + 1
		for piece := engine.PAWN; piece <= engine.KING; piece++ {
			material := points[piece]
			bb := p.PieceBB[color][piece]
			for bb != 0 {
				sq := engine.PopLSB(&bb)
				score += sign * (material + pst[piece][sq^56*(color^1)])
			}
		}
	}
	return score * (-2*p.Stm + 1)
}
