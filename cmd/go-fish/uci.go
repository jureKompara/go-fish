package main

import (
	"fmt"
	"go-fish/internal/engine"
	"os"
	"strings"
	"time"
)

func handleUci(req string, p *engine.Position) {
	cmd, rest, _ := strings.Cut(req, " ")

	switch cmd {

	case "uci":
		greeting()

	case "isready":
		fmt.Println("readyok")

	case "position":
		handlePosition(rest, p)

	case "go":
		{
			abNodes = 0
			qNodes = 0
			ttCutoffs = 0
			TTHit = 0
			TTProbe = 0

			start := time.Now()
			fmt.Println("bestmove", RootSearch(p, 8).Uci())
			elapsed := time.Since(start).Seconds()

			//fmt.Println("ttCutoffs: ", ttCutoffs)
			//fmt.Println("ttHits: ", TTHit)
			//fmt.Println("ttProbes: ", TTProbe)

			nps := int(float64(abNodes+qNodes) / elapsed)

			fmt.Printf("info nodes %d qnodes %d nps %d\n", abNodes, qNodes, nps)
			return
		}
	case "d":
		display(p)

	case "quit":
		os.Exit(0)

	case "stop":
		os.Exit(0)
	}
}

// this expects:
// startpos [moves m1 m2....]
// fen fensStr [moves m1 m2....]
func handlePosition(line string, p *engine.Position) {
	// tokens after "position"
	tok := strings.Fields(line)
	if len(tok) == 0 {
		return
	}

	var off int
	switch tok[0] {
	case "startpos":
		*p = engine.StartPos()
		off = 1

	case "fen":
		// UCI fen is exactly 6 fields: piece placement, stm, castling, ep, halfmove, fullmove
		if len(tok) <= 6 {
			return
		}
		off = 7
		fen := strings.Join(tok[1:7], " ")
		*p = engine.FromFen(fen)

	default:
		return
	}

	if off < len(tok) && tok[off] == "moves" {
		playMoves(p, tok[off+1:])
	}
}

func playMoves(p *engine.Position, tokens []string) {

	for _, t := range tokens {
		moves := p.Movebuff[p.Ply][:]
		n := p.GenMoves(moves)
		moves = moves[:n]

		found := false

		for _, m := range moves {
			if m.Uci() == t {
				found = true
				p.Make(m)
				break
			}
		}
		if !found {
			fmt.Printf("[Err: %v is not legal]\n", t)

		}
	}
}

func greeting() {
	fmt.Println("id name go-fish")
	fmt.Println("id author J")
	fmt.Println("uciok")
}

func display(p *engine.Position) {
	for row := 7; row >= 0; row-- {
		fmt.Println("+---+---+---+---+---+---+---+---+")
		for file := range 8 {
			sq := row*8 + file
			black := uint8(0)
			if p.ColorOcc[engine.BLACK]&(1<<sq) != 0 {
				black = 32
			}

			c := engine.PieceToChar[p.Board[sq]] + black

			fmt.Printf("| %c ", c)
		}
		fmt.Printf("| %d\n", row+1)
	}
	fmt.Println("+---+---+---+---+---+---+---+---+")
	fmt.Println("  a   b   c   d   e   f   g   h")

	fmt.Println("")
	fmt.Println("Fen:", p.ExportFen())
	fmt.Printf("Key: %X\n", p.Hash)
}
