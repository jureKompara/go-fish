package engine

const TTSize = 22

var TT [1 << TTSize]TTEntry
var IndexMask = (uint64(1) << TTSize) - 1

type TTEntry struct {
	Hash      uint64
	Depth     uint8
	Score     int
	BoundType uint8
	BestMove  Move
}
