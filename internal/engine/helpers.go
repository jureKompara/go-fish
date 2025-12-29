package engine

import (
	"fmt"
	"math/bits"
)

var between [64][64]uint64
var line [64][64]uint64

func set(bbptr *uint64, index int) {
	*bbptr |= 1 << index
}

func has(b uint64, sq int) bool {
	return (b & (1 << sq)) != 0
}

func PopLSB(bb *uint64) int {
	lsb_ix := bits.TrailingZeros64(*bb)
	*bb &= *bb - 1
	return lsb_ix
}

// prints a bitboard as a chessboard
func PrintBB(bb uint64) {
	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			if has(bb, rank*8+file) {
				fmt.Print(" X ")
			} else {
				fmt.Print(" _ ")
			}
		}
		fmt.Print("\n")
	}
}

func Between(sq1, sq2 int) uint64 {
	if sq1 == sq2 {
		return uint64(0)
	}

	out := uint64(0)
	r1 := rank(sq1)
	f1 := file(sq1)

	r2 := rank(sq2)
	f2 := file(sq2)

	dr := r2 - r1
	df := f2 - f1

	//if they are not on the same diag return zeroes
	if dr != 0 && df != 0 && !(dr == df || dr == -df) {
		return uint64(0)
	}

	if dr > 0 {
		dr = 1
	}
	if df > 0 {
		df = 1
	}
	if dr < 0 {
		dr = -1
	}
	if df < 0 {
		df = -1
	}

	sq1 += df + dr*8
	for sq1 != sq2 {
		out |= 1 << sq1
		sq1 += df + dr*8
	}

	return out
}

func Line(sq1, sq2 int) uint64 {
	if sq1 == sq2 {
		return uint64(0)
	}

	out := uint64(0)
	r1 := rank(sq1)
	f1 := file(sq1)

	r2 := rank(sq2)
	f2 := file(sq2)

	dr := r2 - r1
	df := f2 - f1

	//if they are not on the same diag return zeroes
	if dr != 0 && df != 0 && !(dr == df || dr == -df) {
		return uint64(0)
	}

	if dr > 0 {
		dr = 1
	}
	if df > 0 {
		df = 1
	}
	if dr < 0 {
		dr = -1
	}
	if df < 0 {
		df = -1
	}
	sq1 += df + dr*8
	for sq1 != sq2 {
		out |= 1 << sq1
		sq1 += df + dr*8
	}
	out |= 1 << sq1

	return out
}

func rank(sq int) int {
	return sq >> 3
}

func file(sq int) int {
	return sq & 7
}
