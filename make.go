package main

func (p *Position) Make(move Move) {

	flags := move.Flags()

	us := p.Stm
	enemy := us ^ 1

	fr := move.From()
	to := move.To()
	piece := p.Board[fr]
	capture := p.Board[to]
	p.save(capture)

	frMask := uint64(1) << fr
	//our piece will always leave from
	p.PieceBB[us][piece] ^= frMask
	p.ColorOcc[us] ^= frMask
	p.Board[fr] = EMPTY

	toMask := uint64(1) << to
	//our piece will always end up at to
	p.PieceBB[us][piece] ^= toMask
	p.ColorOcc[us] ^= toMask
	p.Board[to] = piece
	//p.Key ^= zobristPiece[us][piece][to]

	if capture == EMPTY {
		p.Occ ^= frMask | toMask // empty-to: toggle both
	} else {
		p.PieceBB[enemy][capture] ^= toMask
		p.ColorOcc[enemy] ^= toMask
		p.Occ ^= frMask // capture: only fr becomes empty
	}

	if piece == KING {
		p.kings[us] = to
		if IsCastle(flags) {
			homeRank := us * 56
			t := homeRank + 5 + int(flags)*-2
			tMask := uint64(1) << t
			p.PieceBB[us][ROOK] ^= tMask
			p.ColorOcc[us] ^= tMask
			p.Board[t] = ROOK

			f := homeRank + 7*(1-int(flags))
			fMask := uint64(1) << f
			//p.Key ^= zobristPiece[us][ROOK][7+homeRank]
			p.PieceBB[us][ROOK] ^= fMask
			p.ColorOcc[us] ^= fMask
			p.Board[f] = EMPTY
			p.Occ ^= fMask | tMask
			//p.Key ^= zobristPiece[us][ROOK][5+homeRank]
		}
	}

	/*if p.epSquare != 64 { // REMOVE old EP
		p.Key ^= zobristEP[p.epSquare&7]
	}*/
	p.epSquare = 64

	if flags != QUIET && flags != CAPTURE {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting epSquare
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8 + 16*us
		switch {
		case IsEP(flags):
			epMask := uint64(1) << ep
			p.PieceBB[enemy][PAWN] ^= epMask
			p.ColorOcc[enemy] ^= epMask
			p.Occ ^= epMask
			p.Board[ep] = EMPTY
			//p.Key ^= zobristPiece[enemy][PAWN][ep]

		case IsPromo(flags):
			promo := Promo(flags)
			p.PieceBB[us][PAWN] ^= toMask
			p.PieceBB[us][promo] ^= toMask
			p.Board[to] = promo
			//p.Key ^= zobristPiece[us][PAWN][to]
			//p.Key ^= zobristPiece[us][promo][to]
		case IsDP(flags):
			p.epSquare = uint8(ep)
			//p.Key ^= zobristEP[ep&7]
		}
	}

	//p.Key ^= zobristCastle[p.castleRights]
	//castling rights when rooks move or get capped
	//mby this branch is not worth it?
	if p.castleRights != 0 {
		p.castleRights &= castleMask[fr] & castleMask[to]
	}

	//p.Key ^= zobristCastle[p.castleRights]
	p.Stm ^= 1
	//p.Key ^= zobristSide
	p.Ply++
}
