package main

import (
	"fmt"

	"github.com/fatih/color"
)

func Test(maxDepth int) uint64 {
	ok := color.New(color.Bold, color.FgGreen)
	err := color.New(color.Bold, color.FgRed)

	total := uint64(0)
	for _, test := range Tests {
		fmt.Println(test.FEN)
		pos := FromFen(test.FEN)
		for d := 1; d < len(test.result); d++ {
			if d > maxDepth {
				break
			}
			nodes := pos.Bulk(d)
			total += nodes
			res := test.result[d]
			if nodes != res {
				err.Printf("%d: Err-%d Expected: %d\n", d, nodes, res)
			} else {
				ok.Printf("%d: OK-%d\n", d, nodes)
			}
		}
	}
	return total
}

func (p *Position) Perft(depth int) uint64 {
	if depth == 0 {
		return 1
	}
	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]
	nodes := uint64(0)
	for _, move := range moves {
		p.Make(move)
		nodes += p.Perft(depth - 1)
		p.Unmake(move)
	}
	return nodes
}

func (p *Position) Bulk(depth int) uint64 {
	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]
	// base case depth==1 we just count legal moves
	if depth <= 1 {
		return uint64(n)
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
	moves := p.Movebuff[p.Ply][:]
	n := p.GenMoves(moves)
	moves = moves[:n]
	nodes := uint64(0)
	for _, move := range moves {
		snapCR := p.castleRights
		snapEP := p.epSquare
		snapTM := p.Stm
		snapAll := p.ColorOcc
		snapOcc := p.Occ
		snapKings := p.kings
		snapPieces := p.PieceBB
		p.Make(move)
		n := uint64(0)
		n = p.Bulk(depth - 1)
		fmt.Printf("%s: %d\n", move.Uci(), n)
		nodes += n
		p.Unmake(move)
		if p.castleRights != snapCR ||
			p.epSquare != snapEP ||
			p.Stm != snapTM ||
			p.ColorOcc != snapAll ||
			p.Occ != snapOcc ||
			p.kings != snapKings ||
			p.PieceBB != snapPieces {
			panic("state mismatch")
		}
	}
	return nodes
}
