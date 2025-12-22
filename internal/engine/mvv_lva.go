package engine

// victim values indexed by your piece enum (N,B,R,Q,P,K)
var mvv = [6]int32{300, 300, 500, 900, 100, 0}

func MvvLvaScore(p *Position, m Move) int32 {
	flag := m.Flags()

	var victimSq int
	if flag == EP {
		// captured pawn is behind the TO square
		if p.Stm == WHITE {
			victimSq = m.To() - 8
		} else {
			victimSq = m.To() + 8
		}
	} else {
		victimSq = m.To()
	}
	attacker := int(p.Board[m.From()])
	victim := int(p.Board[victimSq])

	// big constant to keep all captures ahead of quiets if you mix them
	return 10*mvv[victim] - mvv[attacker]
}
