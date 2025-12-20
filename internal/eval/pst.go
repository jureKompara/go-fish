package eval

import (
	"go-fish/internal/engine"
	"math/bits"
)

var points = [5]int{310, 320, 500, 900, 100}

// returns the evaluation of the position from side to move POV
// counts material and PST to get the eval
func Pst(p *engine.Position) int {
	score := 0
	for piece := range engine.KING {
		material := points[piece]
		score += material * bits.OnesCount64(p.PieceBB[engine.WHITE][piece])
		score -= material * bits.OnesCount64(p.PieceBB[engine.BLACK][piece])
	}

	return score * (1 - 2*p.Stm)
}
