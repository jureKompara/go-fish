package eval

import "go-fish/internal/engine"

const PawnValue int32 = 100
const knightValue int32 = 300
const BishopValue int32 = 310
const RookValue int32 = 500
const QueenValue int32 = 900
const KingValue int32 = 0

var Points = [6]int32{
	knightValue,
	BishopValue,
	RookValue,
	QueenValue,
	PawnValue,
	KingValue,
}

var PSQ [2][6][64]int32

// returns the evaluation of the position from side to move POV
// counts material and PST to get the eval
func Pst(p *engine.Position) int32 {
	var score int32 = 0

	black := p.PieceBB[engine.BLACK]
	white := p.PieceBB[engine.WHITE]

	//knights
	bb := white[engine.KNIGHT]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.KNIGHT][engine.PopLSB(&bb)]
	}
	bb = black[engine.KNIGHT]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.KNIGHT][engine.PopLSB(&bb)]
	}

	//bishops
	bb = white[engine.BISHOP]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.BISHOP][engine.PopLSB(&bb)]
	}
	bb = black[engine.BISHOP]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.BISHOP][engine.PopLSB(&bb)]
	}

	//rooks
	bb = white[engine.ROOK]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.ROOK][engine.PopLSB(&bb)]
	}
	bb = black[engine.ROOK]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.ROOK][engine.PopLSB(&bb)]
	}

	//queens
	bb = white[engine.QUEEN]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.QUEEN][engine.PopLSB(&bb)]
	}
	bb = black[engine.QUEEN]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.QUEEN][engine.PopLSB(&bb)]
	}

	//pawns
	bb = white[engine.PAWN]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.PAWN][engine.PopLSB(&bb)]
	}
	bb = black[engine.PAWN]
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
