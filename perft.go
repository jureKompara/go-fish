package main

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
	if depth == 0 {
		return 1
	}
	moves := p.Movebuff[p.Ply][:0]
	p.GenMoves(&moves)
	nodes := uint64(0)
	for _, move := range moves {
		p.Make(move)
		nodes += p.Perft(depth - 1)
		p.Unmake(move)
	}
	return nodes
}

func (p *Position) Bulk(depth int) uint64 {
	if depth == 0 {
		return 1
	}
	moves := p.Movebuff[p.Ply][:0]
	p.GenMoves(&moves)
	// base case depth==1 we dont go to depth 0
	// we just look if the moves are legal and count them
	if depth == 1 {
		return uint64(len(moves))
	}
	nodes := uint64(0)
	for _, move := range moves {
		p.Make(move)
		nodes += p.Bulk(depth - 1)
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
	moves := p.Movebuff[p.Ply][:0]
	p.GenMoves(&moves)
	nodes := uint64(0)
	for _, move := range moves {
		snapCR := p.castleRights
		snapEP := p.epSquare
		snapHM := p.halfMove
		snapTM := p.Stm
		snapFM := p.fullMove
		snapAll := p.ColorBB
		snapOcc := p.Occ
		snapKings := p.kings
		snapPieces := p.PieceBB
		p.Make(move)
		n := uint64(0)
		if !p.isAttacked(p.kings[p.Stm^1], p.Stm) {
			n = p.Perft(depth - 1)
			fmt.Printf("%s: %d\n", move.Uci(), n)
		}
		nodes += n
		p.Unmake(move)
		if p.castleRights != snapCR ||
			p.epSquare != snapEP ||
			p.halfMove != snapHM ||
			p.Stm != snapTM ||
			p.fullMove != snapFM ||
			p.ColorBB != snapAll ||
			p.Occ != snapOcc ||
			p.kings != snapKings ||
			p.PieceBB != snapPieces {
			panic("state mismatch")
		}
	}
	return nodes
}
