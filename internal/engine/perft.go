package engine

import (
	"fmt"

	"github.com/fatih/color"
)

func Test(maxDepth int) {
	ok := color.New(color.Bold, color.FgGreen)
	err := color.New(color.Bold, color.FgRed)

	for _, test := range Tests {
		fmt.Println(test.FEN)
		pos := FromFen(test.FEN)
		for d, res := range test.result {
			if d > maxDepth {
				break
			}
			nodes := pos.Perft(d)
			if nodes != res {
				err.Printf("%d: Err-%d Expected: %d\n", d, nodes, res)
			} else {
				ok.Printf("%d: OK-%d\n", d, nodes)
			}
		}
		fmt.Println("--------------------------------")
	}
}

func (p *Position) Perft(depth int) uint64 {
	// base case depth==1 we dont go to depth 0
	// we just look if the moves are legal and count them
	if depth == 1 {
		count := uint64(0)
		for _, move := range p.GenMoves() {
			p.Make(move)
			if !p.InCheck(1 - p.ToMove) {
				count++
			}
			p.Unmake(move)
		}
		return count
	}
	if depth == 0 {
		return 1
	}
	nodes := uint64(0)
	for _, move := range p.GenMoves() {
		p.Make(move)
		if !p.InCheck(1 - p.ToMove) {
			nodes += p.Perft(depth - 1)
		}
		p.Unmake(move)
	}
	return nodes
}

// very usefull function for debuging
// prints out the # of nodes after each of the first plys in the position
func (p *Position) PerftDivide(depth int) uint64 {
	if depth == 0 {
		return 1
	}
	nodes := uint64(0)
	for _, move := range p.GenMoves() {
		snapCR := p.castleRights
		snapEP := p.epSquare
		snapHM := p.halfMove
		snapTM := p.ToMove
		snapFM := p.fullMove
		snapAll := p.ColorBB
		snapOcc := p.Occupancy
		snapKings := p.kings
		snapPieces := p.PieceBB
		p.Make(move)
		n := uint64(0)
		if !p.isAttacked(p.kings[p.ToMove^1], p.ToMove) {
			n = p.Perft(depth - 1)
			fmt.Printf("%s: %d\n", move.Uci(), n)
		}
		nodes += n
		p.Unmake(move)
		if p.castleRights != snapCR || p.epSquare != snapEP || p.halfMove != snapHM ||
			p.ToMove != snapTM || p.fullMove != snapFM ||
			p.ColorBB != snapAll || p.Occupancy != snapOcc || p.kings != snapKings || p.PieceBB != snapPieces {
			panic("state mismatch")
		}
	}
	fmt.Printf("TOTAL NODES: %d", nodes)
	return nodes
}
