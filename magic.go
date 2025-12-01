package main

func GenMasks() uint {

	return 0
}

func GenMaskBish(sq int) uint64 {
	var out uint64
	for _, d := range bishOff {
		sq2 := sq + d
		for !(sq2&7 == 7 || sq2>>3 == 7 || sq2&7 == 0 || sq2>>3 == 0) {
			set(&out, sq2)
			sq2 += d
		}

	}
	return out
}

func GenRookMask(sq int) uint64 {
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
