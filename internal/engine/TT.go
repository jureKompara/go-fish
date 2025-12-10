package engine

const TTSize = 20

var TT [1 << TTSize]TTEntry
var IndexMask = (uint64(1) << TTSize) - 1

type TTEntry struct {
	Key       uint64
	Depth     uint8
	Score     int
	BoundType uint8
	BestMove  Move
}
