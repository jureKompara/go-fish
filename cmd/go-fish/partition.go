package main

import (
	"go-fish/internal/engine"
)

// expects any []Move
// partition captures first
// returns number of captures
func headCaptures(moves []engine.Move) int {
	capCount := 0
	for i := range moves {
		if moves[i].IsCapture() {
			moves[capCount], moves[i] = moves[i], moves[capCount]
			capCount++
		}
	}
	return capCount
}

// moves all non captures to the end of the slice
// usefull for slices with mostly captures
// returns the index of the first quiet move in the slice
// wich coresponds to the len of capture moves
func tailQuiets(moves []engine.Move) int {
	quietIdx := len(moves)
	for i := quietIdx - 1; i >= 0; i-- {
		if !moves[i].IsCapture() {
			quietIdx--
			moves[quietIdx], moves[i] = moves[i], moves[quietIdx]
		}
	}
	return quietIdx
}
