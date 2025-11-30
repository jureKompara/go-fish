package main

import "math"

// attack maskes per square for every piece type
var knight [64]uint64
var bishop [64]uint64
var rook [64]uint64
var queen [64]uint64
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

// getPawnBB returns diagonal attack targets for a pawn on `sq`; color=true for white.
func getPawnBB(sq uint8, color uint8) uint64 {
	bb := uint64(1) << sq
	if color == 0 {
		return ((bb << 7) & ^h) | ((bb << 9) & ^a)
	} else {
		return ((bb >> 7) & ^a) | ((bb >> 9) & ^h)
	}
}

func getPawnPushBB(sq uint8, color uint8) uint64 {
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
func getSliderBB(sq uint8, deltas []int8) uint64 {
	var out uint64
	var prevF int8
	var sq2 int8
	for _, d := range deltas {
		sq2 = int8(sq)
		prevF = int8(sq) % 8
		for {
			sq2 += d
			newF := sq2 % 8
			if sq2 > 63 || sq2 < 0 || math.Abs(float64(newF-prevF)) > 1 {
				break
			}
			out |= 1 << sq2
			prevF = newF
		}
	}
	return out
}

// getKingBB is like slider but stops after a single step in each direction.
func getKingBB(sq uint8) uint64 {
	var out uint64
	var prevF int8 = int8(sq) % 8
	var sq2 int8
	for _, d := range queenOff {
		sq2 = int8(sq) + d
		if sq2 > 63 || sq2 < 0 || math.Abs(float64(sq2%8-prevF)) > 1 {
			continue
		}
		out |= 1 << sq2
	}
	return out
}

func getKnightBB(sq uint8) uint64 {
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
	for i := range 64 {
		knight[i] = getKnightBB(uint8(i))
		bishop[i] = getSliderBB(uint8(i), bishOff[:])
		rook[i] = getSliderBB(uint8(i), rookOff[:])
		queen[i] = getSliderBB(uint8(i), queenOff[:])
		king[i] = getKingBB(uint8(i))
		pawn[0][i] = getPawnBB(uint8(i), 0)
		pawn[1][i] = getPawnBB(uint8(i), 1)
		pawnPush[0][i] = getPawnPushBB(uint8(i), 0)
		pawnPush[1][i] = getPawnPushBB(uint8(i), 1)
	}
}
