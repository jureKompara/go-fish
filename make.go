package main

func (p *Position) Make(move Move) {
	fr := move.From()
	to := move.To()

	flags := move.Flags()

	us := p.Stm
	enemy := us ^ 1

	piece := p.PieceAt(fr)
	capture := p.PieceAt(to)
	p.save(capture)

	//our piece will always end up at to
	set(&(p.PieceBB[us][piece]), to)
	set(&(p.ColorBB[us]), to)
	p.Board[to] = piece
	//p.Key ^= zobristPiece[us][piece][to]

	//our piece will always leave from
	clear(&(p.PieceBB[us][piece]), fr)
	clear(&(p.ColorBB[us]), fr)
	p.Board[fr] = EMPTY
	//p.Key ^= zobristPiece[us][piece][fr]

	//if its a capture we remove enemy piece from to
	if IsCapture(flags) && !IsEP(flags) {
		clear(&(p.PieceBB[enemy][capture]), to)
		clear(&(p.ColorBB[enemy]), to)
		//p.Key ^= zobristPiece[enemy][capture][to]
	}

	/*if p.epSquare != 64 { // REMOVE old EP
		p.Key ^= zobristEP[p.epSquare&7]
	}*/
	p.epSquare = 64

	if piece == PAWN {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting epSquare
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8*(1-2*us)
		if IsDP(flags) {
			p.epSquare = uint8(ep)
			//p.Key ^= zobristEP[ep&7]
		} else if IsEP(flags) {
			clear(&(p.PieceBB[enemy][PAWN]), ep)
			clear(&(p.ColorBB[enemy]), ep)
			p.Board[ep] = EMPTY
			//p.Key ^= zobristPiece[enemy][PAWN][ep]
		} else if IsPromo(flags) {
			promo := Promo(flags)
			clear(&(p.PieceBB[us][PAWN]), to)
			set(&(p.PieceBB[us][promo]), to)
			p.Board[to] = promo
			//p.Key ^= zobristPiece[us][PAWN][to]
			//p.Key ^= zobristPiece[us][promo][to]
		}
	}

	//p.Key ^= zobristCastle[p.castleRights]
	//castling rights when rooks move or get capped
	//TODO: mby there is a better way to do this
	if p.castleRights != 0 {
		if to == 7 || fr == 7 {
			p.castleRights &= 0b1110
		}
		if to == 0 || fr == 0 {
			p.castleRights &= 0b1101
		}
		if to == 63 || fr == 63 {
			p.castleRights &= 0b1011
		}
		if to == 56 || fr == 56 {
			p.castleRights &= 0b0111
		}
	}

	if piece == KING {
		p.castleRights &= 0b1100 >> (2 * us)
		p.kings[us] = to
		homeRank := us * 56
		switch flags {
		case KCASTLE:
			clear(&(p.PieceBB[us][ROOK]), 7+homeRank)
			clear(&(p.ColorBB[us]), 7+homeRank)
			p.Board[7+homeRank] = EMPTY
			//p.Key ^= zobristPiece[us][ROOK][7+homeRank]

			set(&(p.PieceBB[us][ROOK]), 5+homeRank)
			set(&(p.ColorBB[us]), 5+homeRank)
			p.Board[5+homeRank] = ROOK
			//p.Key ^= zobristPiece[us][ROOK][5+homeRank]
		case QCASTLE:
			clear(&(p.PieceBB[us][ROOK]), homeRank)
			clear(&(p.ColorBB[us]), homeRank)
			p.Board[homeRank] = EMPTY
			//p.Key ^= zobristPiece[us][ROOK][homeRank]

			set(&(p.PieceBB[us][ROOK]), 3+homeRank)
			set(&(p.ColorBB[us]), 3+homeRank)
			p.Board[3+homeRank] = ROOK
			//p.Key ^= zobristPiece[us][ROOK][3+homeRank]
		}
	}

	p.Occ = p.ColorBB[WHITE] | p.ColorBB[BLACK]

	//p.Key ^= zobristCastle[p.castleRights]
	p.Stm ^= 1
	//p.Key ^= zobristSide
	p.Ply++
}
