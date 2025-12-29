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

	flags := move.Flags()
	isPromo := IsPromo(flags)
	to := move.To()
	piece := p.Board[to]
	if isPromo {
		piece = PAWN
	}

	from := move.From()
	fromMask := uint64(1) << from
	//we always put the piece  on the from square
	p.PieceBB[us][piece] ^= fromMask
	p.ColorOcc[us] ^= fromMask
	p.Board[from] = piece

	toMask := uint64(1) << to
	//we clear the to square
	p.ColorOcc[us] ^= toMask
	p.Board[to] = EMPTY
	p.Occ ^= toMask | fromMask

	if isPromo {
		p.PieceBB[us][Promo(flags)] ^= toMask
	} else {
		p.PieceBB[us][piece] ^= toMask
	}

	if flags == EP {

		behind := to - 8 + 16*us
		behindMask := uint64(1) << behind
		p.PieceBB[enemy][PAWN] ^= behindMask
		p.ColorOcc[enemy] ^= behindMask
		p.Occ ^= behindMask
		p.Board[behind] = PAWN

	} else if IsCapture(flags) {
		capture := state.capture
		p.PieceBB[enemy][capture] ^= toMask
		p.ColorOcc[enemy] ^= toMask
		p.Occ ^= toMask
		p.Board[to] = capture
	}

	if piece == KING {

		p.Kings[us] = from
		if IsCastle(flags) {
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
