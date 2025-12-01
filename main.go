package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/fatih/color"
)

const starting_pos string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {

	GenerateAttackBoards()
	fancy := color.New(color.Bold, color.FgHiMagenta)
	var nodes uint64
	pos := FromFen(starting_pos)

	f, _ := os.Create("cpu.out")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i <= 7; i++ {
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

	/*for i := range 64 {
		pirntBB(GenRookMask(i))
		println("------------------------------")
	}*/

}
