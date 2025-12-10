package engine

const promotionRanks uint64 = 0xFF000000000000FF

// returns all pseudo legal moves in the position
func (p *Position) GenMoves(moves *[]Move) {
	color := p.ToMove

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.PieceBB[p.ToMove][piece]
		for bb != 0 {
			sq := PopLSB(&bb)
			switch piece {
			case PAWN:
				p.genPawnMoves(sq, p.pseudoPawn(sq), moves)
			case KNIGHT:
				p.genGenericMoves(sq, KNIGHT, knight[sq] & ^p.ColorBB[color], moves)
			case BISHOP:
				p.genGenericMoves(sq, BISHOP, p.pseudoBishop(sq) & ^p.ColorBB[color], moves)
			case ROOK:
				p.genGenericMoves(sq, ROOK, p.pseudoRook(sq) & ^p.ColorBB[color], moves)
			case QUEEN:
				p.genGenericMoves(sq, QUEEN, (p.pseudoRook(sq)|p.pseudoBishop(sq)) & ^p.ColorBB[color], moves)
			case KING:
				p.genKingMoves(sq, king[sq] & ^p.ColorBB[color], moves)
			}
		}
	}
}

func (p *Position) GenTactics(moves *[]Move) {
	us := p.ToMove
	them := 1 - us

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.PieceBB[p.ToMove][piece]
		for bb != 0 {
			sq := PopLSB(&bb)
			switch piece {
			case PAWN:
				p.genPawnCaptures(sq, p.pseudoPawnCaptures(sq), moves)
				p.genPromotionPushes(sq, p.pseudoPawnPromos(sq), moves)
			case KNIGHT:
				p.genGenericCaptures(sq, KNIGHT, knight[sq]&p.ColorBB[them], moves)
			case BISHOP:
				p.genGenericCaptures(sq, BISHOP, p.pseudoBishop(sq)&p.ColorBB[them], moves)
			case ROOK:
				p.genGenericCaptures(sq, ROOK, p.pseudoRook(sq)&p.ColorBB[them], moves)
			case QUEEN:
				p.genGenericCaptures(sq, QUEEN, (p.pseudoRook(sq)|p.pseudoBishop(sq))&p.ColorBB[them], moves)
			case KING:
				p.genGenericCaptures(sq, KING, king[sq]&p.ColorBB[them], moves)
			}
		}
	}
}

func (p *Position) genKingMoves(sq int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)
		flags := uint8(0)
		capPiece := int(p.Board[to])
		if capPiece != EMPTY {
			flags |= ISCAP
		}
		*moves = append(*moves, NewMove(sq, to, KING, EMPTY, capPiece, flags))
	}

	us := p.ToMove
	enemy := 1 - us
	homeRank := us * 56
	//queenside castle
	if p.castleRights&(0b0001<<(2*us)) != 0 &&
		!has(p.Occupancy, homeRank+5) &&
		!has(p.Occupancy, homeRank+6) &&
		!p.InCheck(us) &&
		!p.isAttacked(sq+1, enemy) &&
		!p.isAttacked(sq+2, enemy) {
		*moves = append(*moves, NewMove(sq, sq+2, KING, EMPTY, EMPTY, KCASTLE))
	}
	//queenside castle
	if p.castleRights&(0b0010<<(2*us)) != 0 &&
		!has(p.Occupancy, homeRank+3) &&
		!has(p.Occupancy, homeRank+2) &&
		!has(p.Occupancy, homeRank+1) &&
		!p.InCheck(us) &&
		!p.isAttacked(sq-1, enemy) &&
		!p.isAttacked(sq-2, enemy) {
		*moves = append(*moves, NewMove(sq, sq-2, KING, EMPTY, EMPTY, QCASTLE))
	}
}

// returns a bitboard with all posible pawn moves
func (p *Position) pseudoPawn(sq int) uint64 {
	us := p.ToMove
	front := sq + 8 - 16*us
	//if the front square isn't empty
	if has(p.Occupancy, front) {
		return pawn[us][sq] & (p.ColorBB[us^1] | 1<<p.epSquare)
	}
	return (pawnPush[us][sq] & ^p.Occupancy) | pawn[us][sq]&(p.ColorBB[us^1]|1<<p.epSquare)
}

// returns a bitboard with only pawn captures
func (p *Position) pseudoPawnCaptures(sq int) uint64 {
	us := p.ToMove
	return pawn[us][sq] & (p.ColorBB[us^1] | 1<<p.epSquare)
}

// returns a bitboard with all posible pawn promotions
func (p *Position) pseudoPawnPromos(sq int) uint64 {
	us := p.ToMove
	return pawnPush[us][sq] & promotionRanks & ^p.Occupancy
}

// we can feed this function any pawn moves
func (p *Position) genPawnMoves(sq int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)
		flags := uint8(0)

		capPiece := EMPTY

		if diff := to - sq; diff == 16 || diff == -16 {
			flags |= DP
		}

		if to == p.epSquare {
			flags |= EP
			capPiece = PAWN
		} else { //if its not an ep we check if its a capture
			capPiece = int(p.Board[to])
			if capPiece != EMPTY {
				flags |= ISCAP
			}
		}

		if !(to>>3 == 7 || to>>3 == 0) {
			*moves = append(*moves, NewMove(sq, to, PAWN, EMPTY, capPiece, flags))
			continue
		}
		for p := KNIGHT; p <= QUEEN; p++ {
			*moves = append(*moves, NewMove(sq, to, PAWN, p, capPiece, flags))
		}
	}
}

// we can only feed this function pawn captures
func (p *Position) genPawnCaptures(sq int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)
		flags := uint8(0)

		var capPiece int
		if to == p.epSquare {
			flags |= EP
			capPiece = PAWN
		} else {
			flags |= ISCAP
			capPiece = int(p.Board[to])
		}

		if !(to>>3 == 7 || to>>3 == 0) {
			*moves = append(*moves, NewMove(sq, to, PAWN, EMPTY, capPiece, flags))
			continue
		}
		for p := KNIGHT; p <= QUEEN; p++ {
			*moves = append(*moves, NewMove(sq, to, PAWN, p, capPiece, flags))
		}
	}
}

// we can only feed this function pawn captures
func (p *Position) genPromotionPushes(sq int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)

		for p := KNIGHT; p <= QUEEN; p++ {
			*moves = append(*moves, NewMove(sq, to, PAWN, p, EMPTY, uint8(0)))
		}
	}
}

func (p *Position) pseudoRook(sq int) uint64 {
	index := (rookMagics[sq] * (maskR[sq] & p.Occupancy) >> rookShifts[sq])
	return rookAttTable[sq][index]
}

func (p *Position) pseudoBishop(sq int) uint64 {
	index := (bishopMagic[sq] * (maskB[sq] & p.Occupancy) >> bishopShifts[sq])
	return bishopAttTable[sq][index]
}

// generates knight and slider moves becouse they have no special cases
// pawns and kings have promotions and castling so they get their own generators
func (p *Position) genGenericMoves(sq, piece int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)
		flags := uint8(0)
		capPiece := int(p.Board[to])
		if capPiece != EMPTY {
			flags |= ISCAP
		}
		*moves = append(*moves, NewMove(sq, to, piece, EMPTY, capPiece, flags))
	}
}

// Same as genGenericMoves but we only ever feed it captures so we dont need to check
// wether its a captures or not
func (p *Position) genGenericCaptures(sq, piece int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)
		capPiece := int(p.Board[to])
		if capPiece == EMPTY {
			panic("OH NOOOOOOOOOOOOOOOOOOOOOOO!!!")
		}
		*moves = append(*moves, NewMove(sq, to, piece, EMPTY, capPiece, ISCAP))
	}
}
