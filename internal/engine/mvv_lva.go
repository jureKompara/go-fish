package engine

var points = [6]int32{300, 300, 500, 900, 100, 1000}

// victim values indexed by your piece enum (N,B,R,Q,P,K)
// index 6 is for EP
var mvv [6][5]int32

func MvvLvaScore(p *Position, m Move) int32 {
	flags := m.Flags()
	if IsPromo(flags) {
		return mvv[PAWN][p.Board[m.To()]] + 10000
	}

	if flags == EP {
		return 900
	}

	return mvv[p.Board[m.From()]][p.Board[m.To()]]
}

func fillMvv() {
	for attacker := range 6 {
		for victim := range 5 {
			mvv[attacker][victim] = 10*points[victim] - points[attacker]
		}
	}
}
