package eval

import "go-fish/internal/engine"

var Points = [7]int32{310, 320, 500, 900, 100, 0, 100}

var PSQ [2][6][64]int32

// returns the evaluation of the position from side to move POV
// counts material and PST to get the eval
func Pst(p *engine.Position) int32 {
	var score int32 = 0

	//knights
	bb := p.PieceBB[engine.WHITE][engine.KNIGHT]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.KNIGHT][engine.PopLSB(&bb)]
	}
	bb = p.PieceBB[engine.BLACK][engine.KNIGHT]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.KNIGHT][engine.PopLSB(&bb)]
	}

	//bishops
	bb = p.PieceBB[engine.WHITE][engine.BISHOP]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.BISHOP][engine.PopLSB(&bb)]
	}
	bb = p.PieceBB[engine.BLACK][engine.BISHOP]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.BISHOP][engine.PopLSB(&bb)]
	}

	//rooks
	bb = p.PieceBB[engine.WHITE][engine.ROOK]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.ROOK][engine.PopLSB(&bb)]
	}
	bb = p.PieceBB[engine.BLACK][engine.ROOK]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.ROOK][engine.PopLSB(&bb)]
	}

	//queens
	bb = p.PieceBB[engine.WHITE][engine.QUEEN]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.QUEEN][engine.PopLSB(&bb)]
	}
	bb = p.PieceBB[engine.BLACK][engine.QUEEN]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.QUEEN][engine.PopLSB(&bb)]
	}

	//pawns
	bb = p.PieceBB[engine.WHITE][engine.PAWN]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.PAWN][engine.PopLSB(&bb)]
	}
	bb = p.PieceBB[engine.BLACK][engine.PAWN]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.PAWN][engine.PopLSB(&bb)]
	}

	//kings
	score += PSQ[engine.WHITE][engine.KING][p.Kings[engine.WHITE]] -
		PSQ[engine.BLACK][engine.KING][p.Kings[engine.BLACK]]

	return score * int32(1-2*p.Stm)
}

func init() {
	for piece := range 6 {
		material := Points[piece]
		for sq := range 64 {
			PSQ[engine.WHITE][piece][sq] = material + pst[piece][sq^56]
			PSQ[engine.BLACK][piece][sq] = material + pst[piece][sq]
		}
	}
}
