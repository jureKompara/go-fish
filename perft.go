package main

import (
	"fmt"

	"github.com/fatih/color"
)

type TestCase struct {
	FEN    string
	result []uint64
}

var tests = []TestCase{
	{FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		result: []uint64{1, 20, 42069, 8902, 197281, 4865609, 119060324, 3195901860, 84998978956},
	},
	{FEN: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		result: []uint64{1, 48, 2039, 97862, 4085603, 193690690, 8031647685},
	},
	{FEN: "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		result: []uint64{1, 14, 191, 2812, 43238, 674624, 11030083, 178633661, 3009794393},
	},
	{FEN: "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		result: []uint64{1, 6, 264, 9467, 422333, 15833292, 706045033},
	},
	{FEN: "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		result: []uint64{1, 44, 1486, 62379, 2103487, 89941194},
	},
	{FEN: "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		result: []uint64{1, 46, 2079, 89890, 3894594, 164075551, 6923051137, 287188994746, 11923589843526},
	},
}

func test(upto int) {

	ok := color.New(color.Bold, color.FgGreen)
	err := color.New(color.Bold, color.FgRed)

	for _, test := range tests {
		fmt.Println(test.FEN)
		pos := FromFen(test.FEN)
		for d, res := range test.result {
			if d > upto {
				break
			}
			nodes := pos.perft(d)
			if nodes != res {
				err.Printf("%d: Err-%d Expected: %d\n", d, nodes, res)
			} else {
				ok.Printf("%d: OK-%d\n", d, nodes)
			}
		}
		fmt.Println("--------------------------------")
	}
}

func (p *Position) perft(depth int) uint64 {
	if depth == 0 {
		return 1
	}
	nodes := uint64(0)
	for _, move := range p.pseudoAll() {
		p.Make(move)
		if !p.IsAttacked(p.kings[1-p.to_move], p.to_move) {
			nodes += p.perft(depth - 1)
		}
		p.Unmake(move)
	}
	return nodes
}

func (p *Position) perftDivide(depth int) uint64 {
	if depth == 0 {
		return 1
	}
	nodes := uint64(0)
	for _, move := range p.pseudoAll() {
		snapCR := p.castle_rights
		snapEP := p.ep_square
		snapHM := p.half_move
		snapTM := p.to_move
		snapFM := p.full_move
		snapAll := p.allBB
		snapOcc := p.occupant
		snapKings := p.kings
		snapPieces := p.pieceBB
		p.Make(move)
		n := uint64(0)
		if !p.IsAttacked(p.kings[1-p.to_move], p.to_move) {
			n = p.perft(depth - 1)
			fmt.Printf("%s: %d\n", move.San(), n)
		}
		nodes += n
		p.Unmake(move)
		if p.castle_rights != snapCR || p.ep_square != snapEP || p.half_move != snapHM ||
			p.to_move != snapTM || p.full_move != snapFM ||
			p.allBB != snapAll || p.occupant != snapOcc || p.kings != snapKings || p.pieceBB != snapPieces {
			panic("state mismatch")
		}
	}
	fmt.Printf("TOTAL NODES: %d", nodes)
	return nodes
}
