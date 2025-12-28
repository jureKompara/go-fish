package main

import (
	"fmt"
	"go-fish/internal/engine"
	"time"
)

// good fens to test
var testFens = []string{
	"7k/5Q2/7K/8/8/8/8/8 w - - 0 1",                                        // mate in 1
	"6k1/1Q6/5K2/8/8/8/8/8 w - - 2 2",                                      // mate in 2
	"r7/2p2p2/p7/2q5/1pP3Pk/1P3p1P/1P2rP2/5RK1 b - c3 0 32",                // mate in 3(king correct)
	"8/4R3/8/5Q2/1P2n3/1B6/PKPk4/8 w - - 7 77",                             //mate in 2(xe4 correct)
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",             //startpos
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", //kiwi
	"r1bqk2r/ppppbppp/2n2n2/4p3/N3P3/3B1N2/PPPP1PPP/R1BQK2R w KQkq - 0 1",  //omar
	"r5k1/5pp1/p1Qpr2p/4p3/4P3/2P2q2/PPP2P1P/4RRK1 w - - 6 21",             //mate in 4(a8 is best!)
	"8/6p1/1P2k1p1/3R2P1/8/5p2/2P2PP1/6K1 w - - 3 48",
	"8/8/2P1Q3/p7/7k/8/P1P4P/6K1 w - - 0 43",
	"r2R2k1/5p1p/5B1P/6P1/p1nP1K2/2P5/8/8 b - - 2 59",
}

func Test(depth, moveTime int) {

	options := Options{
		movetime: moveTime,
		depth:    depth,
	}

	fmt.Println("depth: ", options.depth)

	var fullDuration time.Duration

	for _, fen := range testFens {
		p := engine.FromFen(fen)

		abNodes = 0
		qNodes = 0
		ttCutoffs = 0
		TTHit = 0
		TTProbe = 0

		start := time.Now()
		move := RootSearch(&p, options).Uci()
		elapsed := time.Since(start)
		engine.Killers = [512][2]engine.Move{}
		fullDuration += elapsed

		nps := float64(abNodes+qNodes) / elapsed.Seconds()

		fmt.Println("[", fen, "]")
		fmt.Println("->", move)
		fmt.Println("time: ", elapsed)
		fmt.Printf("info nodes %d qnodes %d @%.2fM nps\n", abNodes, qNodes, nps/1_000_000)

		//fmt.Println("ttProbes:", TTProbe)
		//fmt.Println("TT hit rate:   ", float64(TTHit)/float64(TTProbe)*100, "%")
		//fmt.Println("TT cutoff rate:", float64(ttCutoffs)/float64(TTHit)*100, "%")

		fmt.Println("---------------------------------------------------------------")
	}
	fmt.Println("Test took:", fullDuration)
}
