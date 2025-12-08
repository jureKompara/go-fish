package engine

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

func has(b uint64, sq int) bool {
	return (b & (1 << sq)) != 0
}

func PopLSB(bb *uint64) int {
	lsb_ix := bits.TrailingZeros64(*bb)
	*bb &= *bb - 1
	return lsb_ix
}

// prints a bitboard as a chessboard
func PrintBB(bb uint64, sq int) {
	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			if sq == rank*8+file {
				fmt.Print(" 0 ")
			}
			if has(bb, rank*8+file) {
				fmt.Print(" X ")
			} else {
				fmt.Print(" _ ")
			}
		}
		fmt.Print("\n")
	}
}
