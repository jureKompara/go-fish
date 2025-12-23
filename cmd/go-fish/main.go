package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-fish/internal/engine"
	"os"
	"runtime/pprof"
	"time"
)

var TTProbe = 0
var TTHit = 0
var ttCutoffs = 0

func main() {

	debug := flag.Bool("debug", false, "runs perft on six positions from the wiki")
	prof := flag.Bool("prof", false, "enable profiling")
	perft := flag.Bool("perft", false, "run perft")
	bulk := flag.Bool("bulk", false, "run in bulk mode")
	divide := flag.Bool("divide", false, "run perftDevide")
	test := flag.Bool("test", false, "run test positions for debuging search")
	depth := flag.Int("depth", 5, "set depth for search/perft")
	flag.Parse()

	//enables profiling
	if *prof {

		f, _ := os.Create("cpu.out")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

	}

	if *perft {
		fen := engine.Tests[0].FEN
		pos := engine.FromFen(fen)
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
		nodes := engine.Test(*depth)
		elapsed := time.Since(start)
		fmt.Printf("t: %s\n", elapsed)
		fmt.Printf("N/s: %.2f MN/s\n", float64(nodes)/elapsed.Seconds()/1_000_000)

	} else if *test {

		Test(depth)

	} else if *divide {

		fen := engine.Tests[0].FEN
		pos := engine.FromFen(fen)
		nodes := pos.PerftDivide(*depth)

		fmt.Printf("depth: %d\n", *depth)
		fmt.Printf("FEN: %s\n", fen)
		fmt.Printf("N: %d\n", nodes)

	} else {
		scanner := bufio.NewScanner(os.Stdin)
		p := engine.StartPos()
		for scanner.Scan() {
			handleUci(scanner.Text(), &p)
		}
	}
}
