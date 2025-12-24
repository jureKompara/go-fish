package main

import (
	"fmt"
	"go-fish/internal/engine"
	"time"
)

const INF int32 = 1000000000
const MATE int32 = 1000000

func RootSearch(p *engine.Position, options Options) engine.Move {

	var timeBuget int

	if options.movetime != 0 {
		timeBuget = options.movetime
	} else if p.Stm == engine.WHITE {
		timeBuget = options.wtime/20 + options.winc
	} else {
		timeBuget = options.btime/20 + options.binc
	}
	timeBuget = min(timeBuget, 9000)

	deadline := time.Now().Add(time.Duration(timeBuget) * time.Millisecond)

	moves := p.GenMoves()

	if len(moves) == 0 {
		return 0
	}

	prev := int32(0)
	const base int32 = 25

	for d := 1; d <= options.depth; d++ {

		w := base
		a := prev - w
		b := prev + w

		for {
			bestScore := int32(-INF)
			bestIdx := 0
			alpha := a

			for i, m := range moves {
				// end the search if we are out of time
				if time.Now().After(deadline) {
					fmt.Println("depth ", d-1, "reached")
					return moves[0]
				}

				p.Make(m)

				var score int32
				if i == 0 {
					// 1) full window using aspiration bounds
					score = -AB(p, -b, -alpha, d-1)
				} else {
					// 2) null-window "can this beat alpha?"
					score = -AB(p, -alpha-1, -alpha, d-1)

					// 3) if it might be better, confirm with full window (still within aspiration)
					if score > alpha && score < b {
						score = -AB(p, -b, -alpha, d-1)
					}
				}

				p.Unmake(m)

				if score > bestScore {
					bestScore = score
					bestIdx = i
				}

				if score > alpha {
					alpha = score
				}
			}

			if bestScore <= a { // fail-low: widen
				a -= w
				w *= 2
				continue
			}
			if bestScore >= b { // fail-high: widen
				b += w
				w *= 2
				continue
			}

			prev = bestScore
			moves[0], moves[bestIdx] = moves[bestIdx], moves[0]
			break
		}
	}
	return moves[0]
}
