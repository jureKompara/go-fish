package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	Init()

	debug := flag.Bool("debug", false, "runs perft on six positions from the wiki")
	prof := flag.Bool("prof", false, "enable profiling")
	perft := flag.Bool("perft", false, "run perft")
	bulk := flag.Bool("bulk", false, "run bulk perft")
	divide := flag.Bool("divide", false, "run perftDevide")
	depth := flag.Int("depth", 5, "set depth for search/perft")
	flag.Parse()

	//enables profiling
	if *prof {

		f, _ := os.Create("cpu.out")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

	}

	if *perft {
		fen := Tests[0].FEN
		pos := FromFen(fen)
		start := time.Now()
		var nodes uint64
		if *bulk {
			nodes = pos.Bulk(*depth)
		} else {
			nodes = pos.Perft(*depth)
		}
		elapsed := time.Since(start)
		fmt.Printf("FEN: %s\n", fen)
		fmt.Printf("depth: %d\n", *depth)
		fmt.Printf("t: %s\n", elapsed)
		fmt.Printf("N: %d\n", nodes)
		fmt.Printf("N/s: %.2f MN/s\n", float64(nodes)/elapsed.Seconds()/1_000_000)

	} else if *debug {

		start := time.Now()
		nodes := Test(*depth)
		elapsed := time.Since(start)
		fmt.Printf("t: %s\n", elapsed)
		fmt.Printf("N/s: %.2f MN/s\n", float64(nodes)/elapsed.Seconds()/1_000_000)

	} else if *divide {
		fen := Tests[0].FEN
		pos := FromFen(fen)
		nodes := pos.PerftDivide(*depth)

		fmt.Printf("depth: %d\n", *depth)
		fmt.Printf("FEN: %s\n", fen)
		fmt.Printf("N: %d\n", nodes)
	}
}
