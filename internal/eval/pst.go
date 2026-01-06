package eval

const PawnValue int32 = 100
const knightValue int32 = 310
const BishopValue int32 = 320
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

/*
// [piece][square] from Black's POV
var pst = [6][64]int32{

	// KNIGHT
	{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50,
	},

	// BISHOP
	{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	},
	//ROOK
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		10, 10, 10, 10, 10, 10, 10, 10,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 5, 10, 10, 5, 0, 0,
	},
	//QUEEN
	{
		-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 5, 0, 0, 5, 0, -10,
		-10, 5, 5, 5, 5, 5, 5, -10,
		-5, 0, 5, 5, 5, 5, 0, -5,
		0, 0, 5, 5, 5, 5, 0, 0,
		-10, 5, 5, 5, 5, 5, 5, -10,
		-10, 0, 5, 0, 0, 5, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20,
	},
	// PAWN
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	//KING EARLY/MIDDLE GAME
	{
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20,
	},
	//KING ENDGAME TODO
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

	//knight
	bb := white[engine.KNIGHT]
	for bb != 0 {
		score += PSQ[engine.WHITE][engine.KNIGHT][engine.PopLSB(&bb)]
	}
	bb = black[engine.KNIGHT]
	for bb != 0 {
		score -= PSQ[engine.BLACK][engine.KNIGHT][engine.PopLSB(&bb)]
	}

	//bishop
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
*/
