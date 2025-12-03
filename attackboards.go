package main

//this file contains the
// NOTE:slider pieces are comented out in this file becouse magics makes them irrelevant
// attack maskes per square for every piece type

// var bishop [64]uint64
// var rook [64]uint64
// var queen [64]uint64
var knight [64]uint64
var king [64]uint64
var pawn [2][64]uint64
var pawnPush [2][64]uint64

// file masks used to clip wrap-around moves (e.g., knight moves off the board)
var a uint64
var ab uint64
var h uint64
var gh uint64

func FileMasks() {
	for rank := range 8 {
		a |= (1 << (rank * 8))
		ab |= (3 << (rank * 8))
		h |= (1 << (rank*8 + 7))
		gh |= (3 << (rank*8 + 6))
	}
}

// getPawnBB returns diagonal attack targets for a pawn on `sq`; color=0 for white.
// we clip the a an h files if the pawn is capturing to the right and left respectively
func getPawnBB(sq, color int) uint64 {
	bb := uint64(1) << sq
	if color == 0 {
		return ((bb << 7) & ^h) | ((bb << 9) & ^a)
	} else {
		return ((bb >> 7) & ^a) | ((bb >> 9) & ^h)
	}
}

func getPawnPushBB(sq, color int) uint64 {
	bb := uint64(1) << sq
	var out uint64
	if color == 0 {
		out |= bb << 8
		if 8 <= sq && sq < 16 {
			out |= bb << 16
		}
	} else {
		out |= bb >> 8
		if 48 <= sq && sq < 56 {
			out |= bb >> 16
		}
	}
	return out
}

// getSliderBB walks rays in `deltas` until an edge is hit; returns attack set.
/*func getSliderBB(sq int, deltas []int) uint64 {
	var out uint64
	for _, d := range deltas {
		sq2 := sq
		prevF := sq & 7
		for {
			sq2 += d
			newF := sq2 & 7
			if sq2 > 63 || sq2 < 0 || math.Abs(float64(newF-prevF)) > 1 {
				break
			}
			out |= 1 << sq2
			prevF = newF
		}
	}
	return out
}
*/

// getKingBB is like slider but stops after a single step in each direction.
func getKingBB(sq int) uint64 {
	var out uint64
	prevF := sq & 7
	for _, d := range queenOff {
		sq2 := sq + d
		newF := sq2 & 7
		if sq2 > 63 || sq2 < 0 || newF-prevF > 1 || newF-prevF < -1 {
			continue
		}
		out |= 1 << sq2
	}
	return out
}

// Knight moves on the edge of the board wrap around
// so we have to take out the files that they could wrap to.
// This means taking out the a file for any move going one file to the right
// and taking out both a and b files for moves going two files to the right etc.
func getKnightBB(sq int) uint64 {
	bb := uint64(1) << sq
	return ((bb << 17) & ^a) |
		((bb >> 15) & ^a) |
		((bb << 10) & ^ab) |
		((bb >> 6) & ^ab) |
		((bb << 15) & ^h) |
		((bb >> 17) & ^h) |
		((bb << 6) & ^gh) |
		((bb >> 10) & ^gh)
}

func GenerateAttackBoards() {
	FileMasks()
	for sq := range 64 {
		knight[sq] = getKnightBB(sq)
		//bishop[i] = getSliderBB(i, bishOff[:])
		//rook[i] = getSliderBB(i, rookOff[:])
		//queen[i] = getSliderBB(i, queenOff[:])
		king[sq] = getKingBB(sq)
		pawn[0][sq] = getPawnBB(sq, 0)
		pawn[1][sq] = getPawnBB(sq, 1)
		pawnPush[0][sq] = getPawnPushBB(sq, 0)
		pawnPush[1][sq] = getPawnPushBB(sq, 1)
	}
}
