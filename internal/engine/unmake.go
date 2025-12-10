package engine

func (p *Position) Unmake(move Move) {
	state := p.stateStack[p.Ply-1]
	from := move.From()
	to := move.To()
	piece := move.Piece()
	capture := move.Capture()
	promo := move.Promo()
	flags := move.Flags()

	enemy := p.ToMove
	us := enemy ^ 1

	//we always put the piece  on the from square
	set(&(p.PieceBB[us][piece]), from)
	set(&(p.ColorBB[us]), from)
	set(&(p.Occupancy), from)
	p.Board[from] = uint8(piece)
	p.Key ^= zobristPiece[us][piece][from]

	//we clear the to square
	clear(&(p.PieceBB[us][piece]), to)
	clear(&(p.ColorBB[us]), to)
	clear(&(p.Occupancy), to)
	p.Board[to] = uint8(EMPTY)
	p.Key ^= zobristPiece[us][piece][to]

	if promo != EMPTY {
		clear(&(p.PieceBB[us][promo]), to)
		p.Key ^= zobristPiece[us][PAWN][to]
		p.Key ^= zobristPiece[us][promo][to]

	}

	//if its a cap we put the enemy piece back
	if move.IsCapture() {
		set(&(p.PieceBB[enemy][capture]), to)
		set(&(p.ColorBB[enemy]), to)
		set(&(p.Occupancy), to)
		p.Board[to] = uint8(capture)
		p.Key ^= zobristPiece[enemy][capture][to]
	} else if move.IsEP() {
		behind := to + 8*(1-2*enemy)
		set(&(p.PieceBB[enemy][PAWN]), behind)
		set(&(p.ColorBB[enemy]), behind)
		set(&(p.Occupancy), behind)
		p.Board[behind] = uint8(PAWN)
		p.Key ^= zobristPiece[enemy][PAWN][behind]
	}

	if piece == KING {
		p.kings[us] = from
		homeRank := us * 56
		if flags&KCASTLE != 0 {
			set(&(p.PieceBB[us][ROOK]), 7+homeRank)
			set(&(p.ColorBB[us]), 7+homeRank)
			set(&(p.Occupancy), 7+homeRank)
			p.Board[7+homeRank] = uint8(ROOK)
			p.Key ^= zobristPiece[us][ROOK][7+homeRank]

			clear(&(p.PieceBB[us][ROOK]), 5+homeRank)
			clear(&(p.ColorBB[us]), 5+homeRank)
			clear(&(p.Occupancy), 5+homeRank)
			p.Board[5+homeRank] = uint8(EMPTY)
			p.Key ^= zobristPiece[us][ROOK][5+homeRank]

		} else if flags&QCASTLE != 0 {
			set(&(p.PieceBB[us][ROOK]), homeRank)
			set(&(p.ColorBB[us]), homeRank)
			set(&(p.Occupancy), homeRank)
			p.Board[homeRank] = uint8(ROOK)
			p.Key ^= zobristPiece[us][ROOK][homeRank]

			clear(&(p.PieceBB[us][ROOK]), 3+homeRank)
			clear(&(p.ColorBB[us]), 3+homeRank)
			clear(&(p.Occupancy), 3+homeRank)
			p.Board[3+homeRank] = uint8(EMPTY)
			p.Key ^= zobristPiece[us][ROOK][3+homeRank]
		}
	}

	p.Key ^= zobristCastle[p.castleRights]
	p.castleRights = state.castleRights
	p.Key ^= zobristCastle[p.castleRights]

	if p.epSquare != 64 {
		p.Key ^= zobristEP[p.epSquare&7]
	}

	p.epSquare = int(state.epSquare)
	if p.epSquare != 64 {
		p.Key ^= zobristEP[p.epSquare&7]
	}
	p.halfMove = int(state.halfmove)

	if us == BLACK {
		p.fullMove--
	}
	p.ToMove = us
	p.Key ^= zobristSide
	p.Ply--
}
