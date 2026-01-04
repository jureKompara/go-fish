package engine

const TTSize = 22

type TTBucket struct{ A, B TTEntry }

var TT [1 << TTSize]TTBucket
var IndexMask = uint64(1)<<TTSize - 1

type TTEntry struct {
	Hash      uint64
	Score     int32
	HashMove  Move
	Depth     uint8
	BoundType uint8
}

func (b *TTBucket) Probe(hash uint64) *TTEntry {
	if b.A.Hash == hash {
		return &b.A
	}
	if b.B.Hash == hash {
		return &b.B
	}
	return nil
}

func (b *TTBucket) Store(hit *TTEntry, e TTEntry) {

	if hit != nil {
		*hit = e
		return
	}

	if b.A.Depth <= b.B.Depth {
		b.A = e
	} else {
		b.B = e
	}

}
