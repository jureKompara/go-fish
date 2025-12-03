package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var maskR [64]uint64
var occR [64][4096]uint64
var bishopAttTable [64][4096]uint64

var maskB [64]uint64
var occB [64][4096]uint64
var rookAttTable [64][4096]uint64
var rookShifts [64]int
var bishopShifts [64]int

func main() {
	//init stuff
	GenerateAttackBoards()
	MagicInit()
	pos := FromFen(tests[0].FEN)

	/*profiling stuff
	f, _ := os.Create("cpu.out")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	//pos.perft(7)
	*/

	fancy := color.New(color.Bold, color.FgHiMagenta)
	var nodes uint64
	for i := 1; i <= 6; i++ {
		start := time.Now()
		nodes = pos.perft(i)
		elapsed := time.Since(start)

		seconds := elapsed.Seconds()
		mnps := float64(nodes) / seconds / 1_000_000

		fmt.Printf("depth: %d\n", i)
		fmt.Printf("t: %s\n", elapsed)
		fancy.Printf("N: %d\n", nodes)
		fmt.Printf("N/s: %.2f MN/s\n", mnps)
		fmt.Println("----------------------------------")
	}
	test(5)
}
