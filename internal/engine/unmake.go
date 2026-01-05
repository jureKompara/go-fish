package engine

func (p *Position) Unmake(move Move) {
	p.Ply--
	state := p.stateStack[p.Ply]

	p.Hash = state.hash
	p.castleRights = state.castleRights
	p.epSquare = int(state.epSquare)
	p.HalfMove = int(state.halfmove)

	p.Stm ^= 1
	us := p.Stm
	enemy := us ^ 1

	to := move.To()
	toMask := uint64(1) << to

	var piece uint8 = PAWN
	if !move.IsPromo() {
		piece = p.Board[to]
		p.PieceBB[us][piece] ^= toMask

	} else {
		p.PieceBB[us][move.Promo()] ^= toMask
	}

	from := move.From()
	fromMask := uint64(1) << from
	//we always put the piece  on the from square
	p.PieceBB[us][piece] ^= fromMask
	p.ColorOcc[us] ^= fromMask
	p.Board[from] = piece

	//we clear the to square
	p.ColorOcc[us] ^= toMask
	p.Board[to] = EMPTY
	p.Occ ^= toMask | fromMask

	if move.IsEP() {

		behind := to - 8 + 16*us
		p.Board[behind] = PAWN
		behindMask := uint64(1) << behind
		p.PieceBB[enemy][PAWN] ^= behindMask
		p.ColorOcc[enemy] ^= behindMask
		p.Occ ^= behindMask

	} else if move.IsCapture() {

		capture := state.capture
		p.Board[to] = capture
		p.PieceBB[enemy][capture] ^= toMask
		p.ColorOcc[enemy] ^= toMask
		p.Occ ^= toMask
	}

	if piece == KING {
		p.Kings[us] = from
		if move.IsCastle() {
			flags := move.Flags()
			homeRank := us * 56

			t := homeRank + 5 - 2*int(flags)
			tMask := uint64(1) << t
			p.PieceBB[us][ROOK] ^= tMask
			p.ColorOcc[us] ^= tMask
			p.Board[t] = EMPTY

			f := homeRank + 7*(1-int(flags))
			fMask := uint64(1) << f
			p.PieceBB[us][ROOK] ^= fMask
			p.ColorOcc[us] ^= fMask
			p.Board[f] = ROOK

			p.Occ ^= tMask | fMask
		}
	}
}
