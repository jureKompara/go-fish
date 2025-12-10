package engine

func (p *Position) Make(move Move) {
	p.save()
	to := move.To()
	fr := move.From()
	piece := move.Piece()
	promo := move.Promo()
	flags := move.Flags()
	us := p.ToMove
	enemy := us ^ 1

	//our piece will always end up at to
	set(&(p.PieceBB[us][piece]), to)
	set(&(p.ColorBB[us]), to)
	set(&(p.Occupancy), to)
	p.Board[to] = uint8(piece)
	p.Key ^= zobristPiece[us][piece][to]

	//our piece will always leave from
	clear(&(p.PieceBB[us][piece]), fr)
	clear(&(p.ColorBB[us]), fr)
	clear(&(p.Occupancy), fr)
	p.Board[fr] = uint8(EMPTY)
	p.Key ^= zobristPiece[us][piece][fr]

	//if its a capture we remove enemy piece from to
	if move.IsCapture() {

		//fmt.Println("CAPTURE IS: ", move.Capture())

		clear(&(p.PieceBB[enemy][move.Capture()]), to)
		clear(&(p.ColorBB[enemy]), to)
		p.Key ^= zobristPiece[enemy][move.Capture()][to]
	}

	if us == BLACK {
		p.fullMove++
	}
	if p.epSquare != 64 { // REMOVE old EP
		p.Key ^= zobristEP[p.epSquare&7]
	}
	p.epSquare = 64

	if piece == PAWN {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting epSquare
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8*(1-2*us)
		if move.IsDP() {
			p.epSquare = ep
			p.Key ^= zobristEP[ep&7]
		}
		if move.IsEP() {
			clear(&(p.PieceBB[enemy][PAWN]), ep)
			clear(&(p.ColorBB[enemy]), ep)
			clear(&(p.Occupancy), ep)
			p.Board[ep] = uint8(EMPTY)
			p.Key ^= zobristPiece[enemy][PAWN][ep]
		}
		if promo != EMPTY {
			clear(&(p.PieceBB[us][PAWN]), to)
			set(&(p.PieceBB[us][promo]), to)
			p.Board[to] = uint8(promo)
			p.Key ^= zobristPiece[us][PAWN][to]
			p.Key ^= zobristPiece[us][promo][to]
		}

		p.halfMove = 0
	} else if move.IsCapture() {
		p.halfMove = 0
	} else {
		p.halfMove++
	}

	p.Key ^= zobristCastle[p.castleRights]
	//castling rights when rooks move or get capped
	if to == 0 || fr == 0 {
		p.castleRights &= 0b1101
	}
	if to == 7 || fr == 7 {
		p.castleRights &= 0b1110
	}
	if to == 56 || fr == 56 {
		p.castleRights &= 0b0111
	}
	if to == 63 || fr == 63 {
		p.castleRights &= 0b1011
	}

	if piece == KING {
		p.castleRights &= 0b1100 >> (2 * us)
		p.kings[us] = to
		homeRank := us * 56
		if flags&KCASTLE != 0 {
			clear(&(p.PieceBB[us][ROOK]), 7+homeRank)
			clear(&(p.ColorBB[us]), 7+homeRank)
			clear(&(p.Occupancy), 7+homeRank)
			p.Board[7+homeRank] = uint8(EMPTY)
			p.Key ^= zobristPiece[us][ROOK][7+homeRank]

			set(&(p.PieceBB[us][ROOK]), 5+homeRank)
			set(&(p.ColorBB[us]), 5+homeRank)
			set(&(p.Occupancy), 5+homeRank)
			p.Board[5+homeRank] = uint8(ROOK)
			p.Key ^= zobristPiece[us][ROOK][5+homeRank]

		} else if flags&QCASTLE != 0 {
			clear(&(p.PieceBB[us][ROOK]), homeRank)
			clear(&(p.ColorBB[us]), homeRank)
			clear(&(p.Occupancy), homeRank)
			p.Board[homeRank] = uint8(EMPTY)
			p.Key ^= zobristPiece[us][ROOK][homeRank]

			set(&(p.PieceBB[us][ROOK]), 3+homeRank)
			set(&(p.ColorBB[us]), 3+homeRank)
			set(&(p.Occupancy), 3+homeRank)
			p.Board[3+homeRank] = uint8(ROOK)
			p.Key ^= zobristPiece[us][ROOK][3+homeRank]
		}
	}
	p.Key ^= zobristCastle[p.castleRights]
	p.ToMove = enemy
	p.Key ^= zobristSide
	p.Ply++
}
