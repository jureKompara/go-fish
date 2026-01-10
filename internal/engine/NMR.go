package engine

func (p *Position) MakeNull() {
	p.Stm ^= 1

	p.save(0xFF)
	p.Ply++
	p.HalfMove++

	if p.epSquare != 64 { // REMOVE old EP
		p.Hash ^= zobristEP[p.epSquare&7]
		p.epSquare = 64
	}

	p.Hash ^= zobristSide
}

func (p *Position) UnmakeNull() {
	p.Stm ^= 1
	p.Ply--
	p.HalfMove--

	state := p.stateStack[p.Ply]

	/*assert that halfmove doesnt desync!!!
	if p.HalfMove != int(state.halfmove) {
		panic("ERROR: halfmove desync")
	}
	*/

	p.epSquare = int(state.epSquare)
	p.Hash = state.hash
}
