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
	p.Ply++

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

	if capture == EMPTY {
		p.Occ ^= frMask | toMask // empty-to: toggle both
	} else {
		p.PieceBB[enemy][capture] ^= toMask
		p.ColorOcc[enemy] ^= toMask
		p.Occ ^= frMask // capture: only fr becomes empty
		p.Hash ^= zobristPiece[enemy][capture][to]
	}

	flags := move.Flags()
	if piece == KING {
		p.Kings[us] = to
		if IsCastle(flags) {
			homeRank := us * 56
			t := homeRank + 5 + int(flags)*-2
			tMask := uint64(1) << t
			p.PieceBB[us][ROOK] ^= tMask
			p.ColorOcc[us] ^= tMask
			p.Board[t] = ROOK
			f := homeRank + 7*(1-int(flags))
			fMask := uint64(1) << f
			p.PieceBB[us][ROOK] ^= fMask
			p.ColorOcc[us] ^= fMask
			p.Board[f] = EMPTY
			p.Occ ^= fMask | tMask
			p.Hash ^= zobristPiece[us][ROOK][7+homeRank]
			p.Hash ^= zobristPiece[us][ROOK][5+homeRank]
		}
	}

	if p.epSquare != 64 { // REMOVE old EP
		p.Hash ^= zobristEP[p.epSquare&7]
	}
	p.epSquare = 64

	if flags != QUIET && flags != CAPTURE {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting epSquare
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8 + 16*us
		switch {
		case flags == EP:
			epMask := uint64(1) << ep
			p.PieceBB[enemy][PAWN] ^= epMask
			p.ColorOcc[enemy] ^= epMask
			p.Occ ^= epMask
			p.Board[ep] = EMPTY
			p.Hash ^= zobristPiece[enemy][PAWN][ep]

		case flags == DOUBLE:
			p.epSquare = ep
			p.Hash ^= zobristEP[ep&7]

		case IsPromo(flags):
			promo := Promo(flags)
			p.PieceBB[us][PAWN] ^= toMask
			p.PieceBB[us][promo] ^= toMask
			p.Board[to] = promo
			p.Hash ^= zobristPiece[us][PAWN][to]
			p.Hash ^= zobristPiece[us][promo][to]

		}
	}

	//castling rights when rooks move or get capped
	if p.castleRights != 0 {
		p.Hash ^= zobristCastle[p.castleRights]
		p.castleRights &= castleMask[fr] & castleMask[to]
		p.Hash ^= zobristCastle[p.castleRights]
	}
	p.Hash ^= zobristSide
}
