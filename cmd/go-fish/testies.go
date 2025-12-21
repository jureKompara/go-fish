package main

import (
	"fmt"
	"go-fish/internal/engine"
	"time"
)

// good fens to test
var testFens = []string{
	"7k/5Q2/7K/8/8/8/8/8 w - - 0 1",                         // mate in 1
	"6k1/1Q6/5K2/8/8/8/8/8 w - - 2 2",                       // mate in 2
	"r7/2p2p2/p7/2q5/1pP3Pk/1P3p1P/1P2rP2/5RK1 b - c3 0 32", // mate in 3(king correct)
	"8/4R3/8/5Q2/1P2n3/1B6/PKPk4/8 w - - 7 77",              //mate in 2(xe4 correct)
	engine.STARTPOS,
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", //kiwi
	"r1bqk2r/ppppbppp/2n2n2/4p3/N3P3/3B1N2/PPPP1PPP/R1BQK2R w KQkq - 0 1",  //omar
}

func Test(depth *int) {
	fmt.Println("depth: ", *depth)
	for _, fen := range testFens {
		p := engine.FromFen(fen)

		abNodes = 0
		qNodes = 0
		ttCutoffs = 0
		TTHit = 0
		TTProbe = 0

		start := time.Now()
		move := RootSearch(&p, *depth).Uci()
		elapsed := time.Since(start)

		nps := float64(abNodes+qNodes) / elapsed.Seconds()

		fmt.Println(fen)
		fmt.Println("->", move)
		fmt.Println("time: ", elapsed)
		fmt.Printf("info nodes %d qnodes %d @%.2fM nps\n", abNodes, qNodes, nps/1_000_000)

		fmt.Println("ttProbes:", TTProbe)
		fmt.Println("ttHits:", TTHit)
		fmt.Println("ttCutoffs:", ttCutoffs)

		fmt.Println()
	}
}
