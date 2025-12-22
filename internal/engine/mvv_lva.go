package engine

// victim values indexed by your piece enum (N,B,R,Q,P,K)
// index 6 is for EP
var mvv = [7]int32{300, 300, 500, 900, 100, 0, 0}

func MvvLvaScore(p *Position, m Move) int32 {

	score := 10*mvv[p.Board[m.To()]] - mvv[p.Board[m.From()]]

	if IsPromo(m.Flags()) {
		return 100000
	}

	//victim*10-attacker
	return score
}
