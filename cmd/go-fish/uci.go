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
		{
			fmt.Println("id name go-fish")
			fmt.Println("id author J")
			fmt.Println("uciok")
		}
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

			fmt.Println("bestmove", RootSearch(p, 7).Uci())

			//fmt.Println("ttCutoffs: ", ttCutoffs)
			//fmt.Println("ttHits: ", TTHit)
			//fmt.Println("ttProbes: ", TTProbe)

			elapsed := time.Since(start).Seconds()
			nps := int(float64(abNodes+qNodes) / elapsed)

			fmt.Printf("info nodes %d qnodes %d nps %d\n", abNodes, qNodes, nps)
			return
		}
	case "quit":
		os.Exit(0)
	case "stop":
		os.Exit(0)
	}
}

func handlePosition(stuff string, p *engine.Position) {
	tokens := strings.Split(stuff, " ")
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case "startpos":
			*p = engine.StartPos()
		case "fen":
			i++ //this needs a fix
			*p = engine.FromFen(tokens[i])
		case "moves":
			i++
			for i < len(tokens) {
				moves := p.Movebuff[p.Ply][:]
				n := p.GenMoves(moves)
				moves = moves[:n]
				from, to, promo := parseUci(tokens[i])

				for _, l := range moves {
					if l.From() == from && l.To() == to &&
						(promo == engine.EMPTY || promo == int(engine.Promo(l.Flags()))) {
						p.Make(l)
						break
					}
				}
				i++
			}
		}
	}
}

// converts UCI notation e4e5 to square space
func parseUci(uci string) (f, t, p int) {
	if len(uci) < 4 {
		return -1, -1, engine.EMPTY
	}
	from := uci[0:2]
	to := uci[2:4]
	promo := engine.EMPTY
	if len(uci) == 5 {
		promo = engine.CharToPiece[uci[4]]
	}
	fr := int(from[0] - 'a' + (from[1]-'1')*8)
	too := int(to[0] - 'a' + (to[1]-'1')*8)

	return fr, too, promo
}
