package main

import (
	"fmt"
	"math/bits"
)

func set(bbptr *uint64, index int) {
	*bbptr |= 1 << index
}

func clear(bbptr *uint64, index int) {
	*bbptr &= ^(1 << index)
}

func Has(b uint64, sq int) bool {
	return (b & (1 << sq)) != 0
}

func popLSB(bb *uint64) int {
	lsb_ix := bits.TrailingZeros64(*bb)
	*bb &= *bb - 1
	return lsb_ix
}

// prints a bitboard as a chessboard
func printBB(bb uint64) {
	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			if Has(bb, rank*8+file) {
				fmt.Print(" X ")
			} else {
				fmt.Print(" _ ")
			}
		}
		fmt.Print("\n")
	}
}
