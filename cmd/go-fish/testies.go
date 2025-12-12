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
	"7k/5Q2/6K1/8/8/8/8/8 b - - 0 1",                        // stalemate
	"r7/2p2p2/p7/2q5/1pP3Pk/1P3p1P/1P2rP2/5RK1 b - c3 0 32", // mate in 3(king correct)
	"8/4R3/8/5Q2/1P2n3/1B6/PKPk4/8 w - - 7 77",              //mate in 2(xe4 correct)
}

func Test(depth *int) {
	for _, fen := range testFens {
		//we create the board state
		p := engine.FromFen(fen)

		abNodes = 0
		qNodes = 0
		start := time.Now()
		// we see what move the engine wants to make
		move := RootSearch(&p, *depth).Uci()
		elapsed := time.Since(start).Seconds()
		nps := int(float64(abNodes+qNodes) / elapsed)
		//prints diagnostic info about the search
		fmt.Println(fen, "->", move)
		fmt.Printf("info nodes %d qnodes %d\n", abNodes, qNodes)
		fmt.Printf("info nodes %d qnodes %d nps %d\n", abNodes, qNodes, nps)
		fmt.Println("------------------------------------")
	}
}
