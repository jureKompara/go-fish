package engine

func (p *Position) Make(move Move) {
	us := p.Stm
	enemy := us ^ 1
	p.Stm ^= 1

	fr := move.From()
	to := move.To()
	piece := p.Board[fr]
	capture := p.Board[to]
	p.save(capture)
	p.HashHistory[p.Ply] = p.Hash
	p.Ply++

	if piece == PAWN || move.IsCapture() {
		p.HalfMove = 0
	} else {
		p.HalfMove++
	}

	frMask := uint64(1) << fr
	//our piece will always leave from
	p.PieceBB[us][piece] ^= frMask
	p.ColorOcc[us] ^= frMask
	p.Board[fr] = EMPTY
	p.Hash ^= zobristPiece[us][piece][fr]

	toMask := uint64(1) << to
	//our piece will always end up at to
	p.PieceBB[us][piece] ^= toMask
	p.ColorOcc[us] ^= toMask
	p.Board[to] = piece
	p.Hash ^= zobristPiece[us][piece][to]

	if capture == EMPTY { //not a capture(could be EP)
		p.Occ ^= frMask | toMask // empty-to: toggle both fr and to
	} else { //capture
		p.PieceBB[enemy][capture] ^= toMask
		p.ColorOcc[enemy] ^= toMask
		p.Occ ^= frMask // capture: only toggle fr
		p.Hash ^= zobristPiece[enemy][capture][to]
	}

	if piece == KING {
		p.Kings[us] = to
		if move.IsCastle() {
			flags := move.Flags()

			homeRank := us * 56

			t := homeRank + 5 - 2*int(flags)
			tMask := uint64(1) << t
			p.PieceBB[us][ROOK] ^= tMask
			p.ColorOcc[us] ^= tMask
			p.Board[t] = ROOK
			p.Hash ^= zobristPiece[us][ROOK][t]

			f := homeRank + 7*(1-int(flags))
			fMask := uint64(1) << f
			p.PieceBB[us][ROOK] ^= fMask
			p.ColorOcc[us] ^= fMask
			p.Board[f] = EMPTY
			p.Hash ^= zobristPiece[us][ROOK][f]

			p.Occ ^= fMask | tMask
		}
	}

	if p.epSquare != 64 { // REMOVE old EP
		p.Hash ^= zobristEP[p.epSquare&7]
		p.epSquare = 64
	}

	//for both double pushes and ep captures the relevant squre
	//is the one behind 'to' so we can use it for setting epSquare
	//after a double push or clearing the pawn after an ep capture
	switch {
	case move.IsEP():
		ep := to - 8 + 16*us
		epMask := uint64(1) << ep
		p.PieceBB[enemy][PAWN] ^= epMask
		p.ColorOcc[enemy] ^= epMask
		p.Occ ^= epMask
		p.Board[ep] = EMPTY
		p.Hash ^= zobristPiece[enemy][PAWN][ep]

	case move.IsDouble():
		ep := to - 8 + 16*us
		p.epSquare = ep
		p.Hash ^= zobristEP[ep&7]

	case move.IsPromo():
		promo := Promo(move.Flags())
		p.PieceBB[us][PAWN] ^= toMask
		p.PieceBB[us][promo] ^= toMask
		p.Board[to] = promo
		p.Hash ^= zobristPiece[us][PAWN][to]
		p.Hash ^= zobristPiece[us][promo][to]
	}

	//castling rights when rooks move or get capped
	if p.castleRights != 0 {
		p.Hash ^= zobristCastle[p.castleRights]
		p.castleRights &= castleMask[fr] & castleMask[to]
		p.Hash ^= zobristCastle[p.castleRights]
	}
	p.Hash ^= zobristSide
}
