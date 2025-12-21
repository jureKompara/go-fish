package engine

import (
	"math/bits"
	"math/rand/v2"
)

var maskR [64]uint64
var occR [64][4096]uint64
var bishopAttTable [64][4096]uint64

var maskB [64]uint64
var occB [64][4096]uint64
var rookAttTable [64][4096]uint64
var rookShifts [64]int
var bishopShifts [64]int

// movement offsets for sliders
var bishOff = [4]int{7, 9, -7, -9}
var rookOff = [4]int{1, 8, -1, -8}

var castleMask [64]uint8

func genBishopMask(sq int) uint64 {
	var out uint64
	for _, d := range bishOff {
		sq2 := sq + d
		if sq2 < 0 {
			continue
		}
		for !(sq2&7 == 7 || sq2>>3 == 7 || sq2&7 == 0 || sq2>>3 == 0) {
			set(&out, sq2)
			sq2 += d
		}
	}
	return out
}

func genRookMask(sq int) uint64 {
	var out uint64
	file := sq & 7
	rank := sq >> 3

	for r := rank + 1; r < 7; r++ {
		set(&out, r*8+file)
	}
	for r := rank - 1; r > 0; r-- {
		set(&out, r*8+file)
	}
	for f := file + 1; f < 7; f++ {
		set(&out, rank*8+f)
	}
	for f := file - 1; f > 0; f-- {
		set(&out, rank*8+f)
	}
	return out
}

// this is a little hacky but it works
// this aproach breaks on semy legal fens
func genCastleMask() {
	for sq := range 64 {
		switch sq {
		case 7:
			castleMask[sq] = 0b1110
		case 0:
			castleMask[sq] = 0b1101
		case 63:
			castleMask[sq] = 0b1011
		case 56:
			castleMask[sq] = 0b0111
		case 4:
			castleMask[sq] = 0b1100
		case 60:
			castleMask[sq] = 0b0011
		default:
			castleMask[sq] = 0b1111
		}
	}
}

// This function inits stuff that is required for magics
// that is the rook and bishop masks and the apropriate
// attack bitboard for every square and relevant occupany possible
// also generates random zobrist keys
func init() {
	for i := range 64 {
		maskR[i] = genRookMask(i)
		maskB[i] = genBishopMask(i)
	}
	genCastleMask()

	//fills the line and between lookup table for pins
	for sq1 := range 64 {
		for sq2 := range 64 {
			between[sq1][sq2] = Between(sq1, sq2)
			line[sq1][sq2] = Line(sq1, sq2)
		}
	}

	//The next two for loops generate all possible relevant occupancies
	//for a square
	for sq, bb := range maskR {
		count := bits.OnesCount64(bb)
		l := make([]int, 0, count)
		for bb != 0 {
			l = append(l, int(PopLSB(&bb)))
		}
		for i := range 1 << count {
			cock := uint64(i)
			for cock != 0 {
				occR[sq][i] |= 1 << l[PopLSB(&cock)]
			}
		}
	}
	for sq, bb := range maskB {
		count := bits.OnesCount64(bb)
		l := make([]int, 0, count)
		for bb != 0 {
			l = append(l, int(PopLSB(&bb)))
		}
		for i := range 1 << count {
			cock := uint64(i)
			for cock != 0 {
				occB[sq][i] |= 1 << l[PopLSB(&cock)]
			}
		}
	}

	//fils rookAttTable with the correct attacks
	for sq := range 64 {
		relBits := bits.OnesCount64(maskR[sq])
		occCount := 1 << relBits
		shift := 64 - relBits
		rookShifts[sq] = shift

		for i := range occCount {
			occ := occR[sq][i]
			idx := (occ * rookMagics[sq]) >> shift
			rookAttTable[sq][idx] = sliderAttacks(sq, occ, rookOff[:])
		}
	}

	//fils bishopAttTable with the correct attacks
	for sq := range 64 {
		relBits := bits.OnesCount64(maskB[sq])
		occCount := 1 << relBits
		shift := 64 - relBits
		bishopShifts[sq] = shift

		for i := range occCount {
			occ := occB[sq][i]
			idx := (occ * bishopMagic[sq]) >> shift
			bishopAttTable[sq][idx] = sliderAttacks(sq, occ, bishOff[:])
		}
	}
	//zobrsit numbers init
	//rand.Seed(1)
	var rng = rand.New(rand.NewPCG(0xdeadbeef, 0x12345678))

	for color := range 2 {
		for piece := 0; piece <= KING; piece++ {
			for sq := range 64 {
				zobristPiece[color][piece][sq] = rng.Uint64()
			}
		}
	}
	zobristSide = rng.Uint64()

	for i := range 16 {
		zobristCastle[i] = rng.Uint64()
	}
	for i := range 8 {
		zobristEP[i] = rng.Uint64()
	}
}

func sliderAttacks(sq int, occ uint64, deltas []int) uint64 {
	var out uint64
	var prevF int
	for _, d := range deltas {
		sq2 := sq
		prevF = sq & 0b111
		for {
			sq2 += d
			if sq2 > 63 || sq2 < 0 {
				break
			}
			newF := sq2 & 0b111
			df := newF - prevF
			if df > 1 || df < -1 {
				break
			}

			out |= 1 << sq2
			prevF = newF

			if has(occ, sq2) {
				break
			}
		}
	}
	return out
}
