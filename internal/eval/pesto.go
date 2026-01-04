package eval

/*
import (
	"go-fish/internal/engine"
	"math/bits"
)

const phaseMax int32 = 24

// PeSTO material values are defined per "PeSTO piece type order":
//
//	4 pawn, 0 knight, 1 bishop, 2 rook, 3 queen, 5 king
var mgValuePesto = [6]int32{337, 365, 477, 1025, 82, 0}
var egValuePesto = [6]int32{281, 297, 512, 936, 94, 0}

// Game phase increments per PeSTO piece type order (pawn..king), max total = 24.
var phaseIncPesto = [6]int32{
	1, // knight
	1, // bishop
	2, // rook
	4, // queen
	0, // pawn
	0, // king
}

// PSQMG[color][enginePiece][sq] and PSQEG[color][enginePiece][sq]
// include (material + PST) already.
var PSQMG [2][6][64]int32
var PSQEG [2][6][64]int32

// --- Black-POV PeSTO PSTs (rank 8 first, rank 1 last) ---
// Order: pawn, knight, bishop, rook, queen, king  (PeSTO order)

var mgPSTBlack = [6][64]int32{

	// knight
	{
		-105, -21, -58, -33, -17, -28, -19, -23,
		-29, -53, -12, -3, -1, 18, -14, -19,
		-23, -9, 12, 10, 19, 17, 25, -16,
		-13, 4, 16, 13, 28, 19, 21, -8,
		-9, 17, 19, 53, 37, 69, 18, 22,
		-47, 60, 37, 65, 84, 129, 73, 44,
		-73, -41, 72, 36, 23, 62, 7, -17,
		-167, -89, -34, -49, 61, -97, -15, -107,
	},
	// bishop
	{
		-33, -3, -14, -21, -13, -12, -39, -21,
		4, 15, 16, 0, 7, 21, 33, 1,
		0, 15, 15, 15, 14, 27, 18, 10,
		-6, 13, 13, 26, 34, 12, 10, 4,
		-4, 5, 19, 50, 37, 37, 7, -2,
		-16, 37, 43, 40, 35, 50, 37, -2,
		-26, 16, -18, -13, 30, 59, 18, -47,
		-29, 4, -82, -37, -25, -42, 7, -8,
	},
	// rook
	{
		-19, -13, 1, 17, 16, 7, -37, -26,
		-44, -16, -20, -9, -1, 11, -6, -71,
		-45, -25, -16, -17, 3, 0, -5, -33,
		-36, -26, -12, -1, 9, -7, 6, -23,
		-24, -11, 7, 26, 24, 35, -8, -20,
		-5, 19, 26, 36, 17, 45, 61, 16,
		27, 32, 58, 62, 80, 67, 26, 44,
		32, 42, 32, 51, 63, 9, 31, 43,
	},
	// queen
	{
		-1, -18, -9, 10, -15, -25, -31, -50,
		-35, -8, 11, 2, 8, 15, -3, 1,
		-14, 2, -11, -2, -5, 2, 14, 5,
		-9, -26, -9, -10, -2, -4, 3, -3,
		-27, -27, -16, -16, -1, 17, -2, 1,
		-13, -17, 7, 8, 29, 56, 47, 57,
		-24, -39, -5, 1, -16, 57, 28, 54,
		-28, 0, 29, 12, 59, 44, 43, 45,
	},
	// pawn
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		-35, -1, -20, -23, -15, 24, 38, -22,
		-26, -4, -4, -10, 3, 3, 33, -12,
		-27, -2, -5, 12, 17, 6, 10, -25,
		-14, 13, 6, 21, 23, 12, 17, -23,
		-6, 7, 26, 31, 65, 56, 25, -20,
		98, 134, 61, 95, 68, 126, 34, -11,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	// king
	{
		-15, 36, 12, -54, 8, -28, 24, 14,
		1, 7, -8, -64, -43, -16, 9, 8,
		-14, -14, -22, -46, -44, -30, -15, -27,
		-49, -1, -27, -39, -46, -44, -33, -51,
		-17, -20, -12, -27, -30, -25, -14, -36,
		-9, 24, 2, -16, -20, 6, 22, -22,
		29, -1, -20, -7, -8, -4, -38, -29,
		-65, 23, 16, -15, -56, -34, 2, 13,
	},
}

var egPSTBlack = [6][64]int32{

	// knight
	{
		-29, -51, -23, -15, -22, -18, -50, -64,
		-42, -20, -10, -5, -2, -20, -23, -44,
		-23, -3, -1, 15, 10, -3, -20, -22,
		-18, -6, 16, 25, 16, 17, 4, -18,
		-17, 3, 22, 22, 22, 11, 8, -18,
		-24, -20, 10, 9, -1, -9, -19, -41,
		-25, -8, -25, -2, -9, -25, -24, -52,
		-58, -38, -13, -28, -31, -27, -63, -99,
	},
	// bishop
	{
		-23, -9, -23, -5, -9, -16, -5, -17,
		-14, -18, -7, -1, 4, -9, -15, -27,
		-12, -3, 8, 10, 13, 3, -7, -15,
		-6, 3, 13, 19, 7, 10, -3, -9,
		-3, 9, 12, 9, 14, 10, 3, 2,
		2, -8, 0, -1, -2, 6, 0, 4,
		-8, -4, 7, -12, -3, -13, -4, -14,
		-14, -21, -11, -8, -7, -9, -17, -24,
	},
	// rook
	{
		-9, 2, 3, -1, -5, -13, 4, -20,
		-6, -6, 0, 2, -9, -9, -11, -3,
		-4, 0, -5, -1, -7, -12, -8, -16,
		3, 5, 8, 4, -5, -6, -8, -11,
		4, 3, 13, 1, 2, 1, -1, 2,
		7, 7, 7, 5, 4, -3, -5, -3,
		11, 13, 13, 11, -3, 3, 8, 3,
		13, 10, 18, 15, 12, 12, 8, 5,
	},
	// queen
	{
		-33, -28, -22, -43, -5, -32, -20, -41,
		-22, -23, -30, -16, -16, -23, -36, -32,
		-16, -27, 15, 6, 9, 17, 10, 5,
		-18, 28, 19, 47, 31, 34, 39, 23,
		3, 22, 24, 45, 57, 40, 57, 36,
		-20, 6, 9, 49, 47, 35, 19, 9,
		-17, 20, 32, 41, 58, 25, 30, 0,
		-9, 22, 22, 27, 27, 19, 10, 20,
	},
	// pawn
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		13, 8, 8, 10, 13, 0, 2, -7,
		4, 7, -6, 1, 0, -5, -1, -8,
		13, 9, -3, -7, -7, -8, 3, -1,
		32, 24, 13, 5, -2, 4, 17, 17,
		94, 100, 85, 67, 56, 53, 82, 84,
		178, 173, 158, 134, 147, 132, 165, 187,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	// king
	{
		-53, -34, -21, -11, -28, -14, -24, -43,
		-27, -11, 4, 13, 14, 4, -5, -17,
		-19, -3, 11, 21, 23, 16, 7, -9,
		-18, -4, 21, 24, 27, 23, 9, -11,
		-8, 22, 24, 27, 26, 33, 26, 3,
		10, 17, 23, 15, 20, 45, 44, 13,
		-12, 17, 14, 17, 17, 38, 23, 11,
		-74, -35, -18, -18, -11, 15, 4, -17,
	},
}

func init() {
	for piece := range 6 {
		mgMat := mgValuePesto[piece]
		egMat := egValuePesto[piece]

		for sq := range 64 {
			//pesto PSQ
			PSQMG[engine.WHITE][piece][sq] = mgMat + mgPSTBlack[piece][sq]
			PSQEG[engine.WHITE][piece][sq] = egMat + egPSTBlack[piece][sq]

			PSQMG[engine.BLACK][piece][sq] = mgMat + mgPSTBlack[piece][sq^56]
			PSQEG[engine.BLACK][piece][sq] = egMat + egPSTBlack[piece][sq^56]
		}
	}
}

// EvalPeSTO returns tapered PeSTO eval from side-to-move POV.
func PeSTO(p *engine.Position) int32 {
	var mgW, mgB int32
	var egW, egB int32
	var phase int32

	white := p.PieceBB[engine.WHITE]
	black := p.PieceBB[engine.BLACK]

	// Accumulate MG/EG PSQ and game phase (counting pieces by type).
	for piece := range 6 {
		// WHITE
		bb := white[piece]
		if bb != 0 {
			phase += int32(bits.OnesCount64(bb)) * phaseIncPesto[piece]
			for bb != 0 {
				sq := engine.PopLSB(&bb)
				mgW += PSQMG[engine.WHITE][piece][sq]
				egW += PSQEG[engine.WHITE][piece][sq]
			}
		}

		// BLACK
		bb = black[piece]
		if bb != 0 {
			phase += int32(bits.OnesCount64(bb)) * phaseIncPesto[piece]
			for bb != 0 {
				sq := engine.PopLSB(&bb)
				mgB += PSQMG[engine.BLACK][piece][sq]
				egB += PSQEG[engine.BLACK][piece][sq]
			}
		}
	}

	if phase > phaseMax {
		phase = phaseMax
	}
	egPhase := phaseMax - phase

	mgScore := mgW - mgB
	egScore := egW - egB

	score := (mgScore*phase + egScore*egPhase) / phaseMax

	// side-to-move POV: WHITE stm -> keep, BLACK stm -> negate
	return score * int32(1-2*p.Stm)
}
*/
