package main

func (p *Position) Unmake(move Move) {

	state := p.stateStack[p.Ply-1]
	from := move.From()
	to := move.To()
	flags := move.Flags()

	enemy := p.Stm
	us := enemy ^ 1

	var piece uint8
	isPromo := IsPromo(flags)

	if isPromo {
		piece = PAWN
	} else {
		piece = p.PieceAt(to)
	}

	//we always put the piece  on the from square
	set(&(p.PieceBB[us][piece]), from)
	set(&(p.ColorBB[us]), from)
	p.Board[from] = piece

	//we clear the to square
	clear(&(p.ColorBB[us]), to)
	p.Board[to] = EMPTY

	if isPromo {
		clear(&(p.PieceBB[us][Promo(flags)]), to)
	} else {
		clear(&(p.PieceBB[us][piece]), to)
	}

	if IsEP(flags) {
		behind := to + 8*(1-2*enemy)
		set(&(p.PieceBB[enemy][PAWN]), behind)
		set(&(p.ColorBB[enemy]), behind)
		p.Board[behind] = PAWN
	} else if IsCapture(flags) {
		capture := state.Capture()
		set(&(p.PieceBB[enemy][capture]), to)
		set(&(p.ColorBB[enemy]), to)
		p.Board[to] = capture
	}

	if piece == KING {
		p.kings[us] = from
		homeRank := us * 56
		switch flags {
		case KCASTLE:
			set(&(p.PieceBB[us][ROOK]), 7+homeRank)
			set(&(p.ColorBB[us]), 7+homeRank)
			p.Board[7+homeRank] = ROOK

			clear(&(p.PieceBB[us][ROOK]), 5+homeRank)
			clear(&(p.ColorBB[us]), 5+homeRank)
			p.Board[5+homeRank] = EMPTY

		case QCASTLE:
			set(&(p.PieceBB[us][ROOK]), homeRank)
			set(&(p.ColorBB[us]), homeRank)
			p.Board[homeRank] = ROOK

			clear(&(p.PieceBB[us][ROOK]), 3+homeRank)
			clear(&(p.ColorBB[us]), 3+homeRank)
			p.Board[3+homeRank] = EMPTY
		}
	}
	p.Occ = p.ColorBB[WHITE] | p.ColorBB[BLACK]

	p.castleRights = state.CastleRights()
	p.epSquare = state.EPsquare()

	p.Stm ^= 1
	p.Ply--
}
