package engine

func (p *Position) Unmake(move Move) {
	state := p.stateStack[p.Ply-1]
	from := move.From()
	to := move.To()
	piece := p.PieceAt(to)
	flags := move.Flags()

	if IsPromo(flags) {
		piece = PAWN
	}

	capture := state.capture

	enemy := p.ToMove
	us := enemy ^ 1

	//we always put the piece  on the from square
	set(&(p.PieceBB[us][piece]), from)
	set(&(p.ColorBB[us]), from)
	set(&(p.Occupancy), from)
	p.Board[from] = piece

	//we clear the to square
	clear(&(p.PieceBB[us][piece]), to)
	clear(&(p.ColorBB[us]), to)
	clear(&(p.Occupancy), to)
	p.Board[to] = EMPTY

	if IsPromo(flags) {
		promo := Promo(flags)
		clear(&(p.PieceBB[us][promo]), to)
	}

	//if its a cap we put the enemy piece back
	if IsEP(flags) {
		behind := to + 8*(1-2*enemy)
		set(&(p.PieceBB[enemy][PAWN]), behind)
		set(&(p.ColorBB[enemy]), behind)
		set(&(p.Occupancy), behind)
		p.Board[behind] = PAWN
	} else if IsCapture(flags) {
		set(&(p.PieceBB[enemy][capture]), to)
		set(&(p.ColorBB[enemy]), to)
		set(&(p.Occupancy), to)
		p.Board[to] = capture
	}

	if piece == KING {
		p.kings[us] = from
		homeRank := us * 56
		switch flags {
		case KCASTLE:
			set(&(p.PieceBB[us][ROOK]), 7+homeRank)
			set(&(p.ColorBB[us]), 7+homeRank)
			set(&(p.Occupancy), 7+homeRank)
			p.Board[7+homeRank] = ROOK

			clear(&(p.PieceBB[us][ROOK]), 5+homeRank)
			clear(&(p.ColorBB[us]), 5+homeRank)
			clear(&(p.Occupancy), 5+homeRank)
			p.Board[5+homeRank] = EMPTY

		case QCASTLE:
			set(&(p.PieceBB[us][ROOK]), homeRank)
			set(&(p.ColorBB[us]), homeRank)
			set(&(p.Occupancy), homeRank)
			p.Board[homeRank] = ROOK

			clear(&(p.PieceBB[us][ROOK]), 3+homeRank)
			clear(&(p.ColorBB[us]), 3+homeRank)
			clear(&(p.Occupancy), 3+homeRank)
			p.Board[3+homeRank] = EMPTY
		}
	}

	p.castleRights = state.castleRights

	p.epSquare = int(state.epSquare)

	p.halfMove = int(state.halfmove)

	if us == BLACK {
		p.fullMove--
	}
	p.ToMove = us
	p.Ply--
}
