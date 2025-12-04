// internal/engine/bench_test.go
package engine

import "testing"

func BenchmarkPerftStartposDepth5(b *testing.B) {
	MagicInit()
	fen := Tests[0].FEN // or just hardcode starting_pos fen
	for b.Loop() {
		pos := FromFen(fen)
		_ = pos.Perft(5)
	}
}
