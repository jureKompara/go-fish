package main

import (
	"fmt"
	"math/bits"
)

func set(bbptr *uint64, index uint8) {
	*bbptr |= 1 << index
}

func clear(bbptr *uint64, index uint8) {
	*bbptr &= ^(1 << index)
}

func Has(b uint64, sq uint8) bool {
	return (b & (1 << sq)) != 0
}

func popLSB(bb *uint64) uint8 {
	lsb_ix := bits.TrailingZeros64(*bb)
	*bb &= *bb - 1
	return uint8(lsb_ix)
}

// prints a bitboard as a chessboard
func pirntBB(bb uint64) {
	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			if Has(bb, uint8(rank*8+file)) {
				fmt.Print(" X ")
			} else {
				fmt.Print(" _ ")
			}
		}
		fmt.Print("\n")
	}
}
