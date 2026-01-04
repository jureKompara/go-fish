package eval

import (
	"go-fish/internal/engine"
)

const PawnValue int32 = 100
const knightValue int32 = 320
const BishopValue int32 = 330
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

func init() {
	for piece := range 6 {
		for sq := range 64 {
			PSQ[engine.WHITE][piece][sq] = pst[piece][sq^56] + Points[piece]
			PSQ[engine.BLACK][piece][sq] = pst[piece][sq] + Points[piece]
		}
	}
}

func Eval(p *engine.Position) int32 {

	var score int32

	white := p.PieceBB[engine.WHITE]
	black := p.PieceBB[engine.BLACK]

	bb := white[engine.KNIGHT]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.KNIGHT][engine.PopLSB(&bb)]
	}
	bb = black[engine.KNIGHT]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.KNIGHT][engine.PopLSB(&bb)]
	}

	bb = white[engine.BISHOP]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.BISHOP][engine.PopLSB(&bb)]
	}
	bb = black[engine.BISHOP]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.BISHOP][engine.PopLSB(&bb)]
	}

	bb = white[engine.ROOK]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.ROOK][engine.PopLSB(&bb)]
	}
	bb = black[engine.ROOK]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.ROOK][engine.PopLSB(&bb)]
	}

	bb = white[engine.QUEEN]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.QUEEN][engine.PopLSB(&bb)]
	}
	bb = black[engine.QUEEN]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.QUEEN][engine.PopLSB(&bb)]
	}

	bb = white[engine.PAWN]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.PAWN][engine.PopLSB(&bb)]
	}
	bb = black[engine.PAWN]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.PAWN][engine.PopLSB(&bb)]
	}

	score += PSQ[engine.WHITE][engine.KING][p.Kings[engine.WHITE]] -
		PSQ[engine.BLACK][engine.KING][p.Kings[engine.BLACK]]

	return score * int32(1-2*p.Stm)
}
