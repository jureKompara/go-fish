package eval

import "go-fish/internal/engine"

var points = [6]int{300, 310, 500, 900, 100, 0}

var PSQ [2][6][64]int

// returns the evaluation of the position from side to move POV
// counts material and PST to get the eval
func Pst(p *engine.Position) int {
	score := 0
	for piece := 0; piece <= engine.KING; piece++ {
		bb := p.PieceBB[engine.WHITE][piece]
		for bb != 0 {
			score += PSQ[engine.WHITE][piece][engine.PopLSB(&bb)]
		}
		bb = p.PieceBB[engine.BLACK][piece]
		for bb != 0 {
			score -= PSQ[engine.BLACK][piece][engine.PopLSB(&bb)]
		}
	}
	return score * (1 - 2*p.Stm)
}

func init() {
	for piece := range 6 {
		material := points[piece]
		for sq := range 64 {
			PSQ[engine.WHITE][piece][sq] = material + pst[piece][sq^56]
			PSQ[engine.BLACK][piece][sq] = material + pst[piece][sq]
		}
	}
}
