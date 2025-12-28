package main

import (
	"go-fish/internal/engine"
	"sort"
)

// expects any []Move
// partitions a move slice and orders it via MVVLVA
// returns number of captures
func partitionSort(p *engine.Position, moves []engine.Move) int {

	//partition captures first
	capCount := 0
	for i := range moves {
		if engine.IsCapture(moves[i].Flags()) {
			moves[capCount], moves[i] = moves[i], moves[capCount]
			capCount++
		}
	}

	//only sort if there is more than 1 capture
	if capCount > 1 {
		sort.Slice(moves[0:capCount], func(i, j int) bool {
			return engine.MvvLvaScore(p, moves[i]) > engine.MvvLvaScore(p, moves[j])
		})
	}
	return capCount
}

// moves all non captures to the end of the slice
// usefull for slices with mostly captures
func tailQuiets(moves []engine.Move) int {
	quietIdx := len(moves)
	for i := quietIdx - 1; i >= 0; i-- {
		if !engine.IsCapture(moves[i].Flags()) {
			quietIdx--
			moves[quietIdx], moves[i] = moves[i], moves[quietIdx]
		}
	}
	return quietIdx
}
