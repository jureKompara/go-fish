package main

import (
	"fmt"
	"go-fish/internal/engine"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleUci(req string, p *engine.Position) {

	cmd, rest, ok := strings.Cut(req, " ")

	switch cmd {

	case "uci":
		greeting()

	case "isready":
		fmt.Println("readyok")

	case "position":
		if ok {
			handlePosition(rest, p)
		}

	case "go":
		{
			options := optionsParser(rest)
			//reset diagnostic values
			abNodes = 0
			qNodes = 0
			TTProbe = 0
			TTHit = 0
			ttCutoffs = 0

			start := time.Now()
			fmt.Println("bestmove", RootSearch(p, options).Uci())
			elapsed := time.Since(start).Seconds()

			nps := int(float64(abNodes+qNodes) / elapsed)

			fmt.Printf("info nodes %d qnodes %d nps %d\n", abNodes, qNodes, nps)
			engine.Killers = [512][2]engine.Move{}
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

// Expects:
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
		if len(tok) < 7 {
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

// handles the playing of moves after a position command
func playMoves(p *engine.Position, tokens []string) {
	for _, t := range tokens {

		found := false
		moves := p.GenMoves()
		for _, m := range moves {
			if m.Uci() == t {
				found = true
				p.Make(m)
				break
			}
		}
		if !found {
			fmt.Printf("[Err: %s is not legal]\n", t)
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

func optionsParser(options string) Options {

	out := Options{
		wtime:    900000000,
		btime:    900000000,
		depth:    64,
		movetime: 0,
		nodes:    0,
		winc:     0,
		binc:     0,
	}

	tokens := strings.Fields(options)

	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case "wtime":
			out.wtime, _ = strconv.Atoi(tokens[i+1])
		case "btime":
			out.btime, _ = strconv.Atoi(tokens[i+1])
		case "movetime":
			out.movetime, _ = strconv.Atoi(tokens[i+1])

		case "winc":
			out.winc, _ = strconv.Atoi(tokens[i+1])
		case "binc":
			out.binc, _ = strconv.Atoi(tokens[i+1])

		case "depth":
			out.depth, _ = strconv.Atoi(tokens[i+1])
		}
		i += 2
	}
	return out
}

type Options struct {
	wtime    int
	btime    int
	depth    int
	movetime int
	nodes    int
	winc     int
	binc     int
}
