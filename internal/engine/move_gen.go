package engine

//used for promo only generation for generating tactical moves
//const promotionRanks uint64 = 0xFF000000000000FF

// returns all pseudo legal moves in the position
func (p *Position) GenMoves(moves *[]Move) {
	color := p.ToMove

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.PieceBB[p.ToMove][piece]
		for bb != 0 {
			sq := PopLSB(&bb)
			switch piece {
			case PAWN:
				p.genPawnMoves(sq, moves)
			case KNIGHT:
				p.genGenericMoves(sq, knight[sq] & ^p.ColorBB[color], moves)
			case BISHOP:
				p.genGenericMoves(sq, p.pseudoBishop(sq) & ^p.ColorBB[color], moves)
			case ROOK:
				p.genGenericMoves(sq, p.pseudoRook(sq) & ^p.ColorBB[color], moves)
			case QUEEN:
				p.genGenericMoves(sq, (p.pseudoRook(sq)|p.pseudoBishop(sq)) & ^p.ColorBB[color], moves)
			case KING:
				p.genKingMoves(sq, moves)
			}
		}
	}
}

func (p *Position) genKingMoves(sq int, moves *[]Move) {

	mask := king[sq] & ^p.ColorBB[p.ToMove]

	captures := mask & p.ColorBB[p.ToMove^1]

	quiets := mask & ^p.Occupancy

	for captures != 0 {
		to := PopLSB(&captures)
		*moves = append(*moves, NewMove(sq, to, CAPTURE))
	}

	for quiets != 0 {
		to := PopLSB(&quiets)
		*moves = append(*moves, NewMove(sq, to, QUIET))
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
		*moves = append(*moves, NewMove(sq, sq+2, KCASTLE))
	}
	//queenside castle
	if p.castleRights&(0b0010<<(2*us)) != 0 &&
		!has(p.Occupancy, homeRank+3) &&
		!has(p.Occupancy, homeRank+2) &&
		!has(p.Occupancy, homeRank+1) &&
		!p.InCheck(us) &&
		!p.isAttacked(sq-1, enemy) &&
		!p.isAttacked(sq-2, enemy) {
		*moves = append(*moves, NewMove(sq, sq-2, QCASTLE))
	}
}

// returns a bitboard with all posible pawn moves
func (p *Position) pseudoPawn(sq int) uint64 {
	us := p.ToMove
	enemy := us ^ 1
	front := sq + 8 - 16*us
	//if the front square is occupied
	if has(p.Occupancy, front) {
		return pawn[us][sq] & (p.ColorBB[enemy] | 1<<p.epSquare)
	}
	return (pawnPush[us][sq] & ^p.Occupancy) | pawn[us][sq]&(p.ColorBB[enemy]|1<<p.epSquare)
}

// we can feed this function any pawn moves
func (p *Position) genPawnMoves(sq int, moves *[]Move) {

	mask := p.pseudoPawn(sq)

	captures := mask & p.ColorBB[p.ToMove^1]

	ep := mask & (1 << p.epSquare)

	quiets := mask & ^(p.Occupancy | ep)

	for captures != 0 {
		to := PopLSB(&captures)
		rank := to >> 3

		if !(rank == 7 || rank == 0) {
			*moves = append(*moves, NewMove(sq, to, CAPTURE))
			continue
		}
		for promoFlag := PROMOQUEENX; promoFlag >= PROMOKNIGHTX; promoFlag-- {
			*moves = append(*moves, NewMove(sq, to, promoFlag))
		}
	}

	for quiets != 0 {
		to := PopLSB(&quiets)

		if diff := to - sq; diff == 16 || diff == -16 {
			*moves = append(*moves, NewMove(sq, to, DOUBLE))
			continue
		} else if !(to>>3 == 7 || to>>3 == 0) {
			*moves = append(*moves, NewMove(sq, to, QUIET))
			continue
		}
		for promoFlag := PROMOQUEEN; promoFlag >= PROMOKNIGHT; promoFlag-- {
			*moves = append(*moves, NewMove(sq, to, promoFlag))
		}
	}

	if ep != 0 {
		to := PopLSB(&ep)
		*moves = append(*moves, NewMove(sq, to, EP))
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
func (p *Position) genGenericMoves(sq int, mask uint64, moves *[]Move) {

	captures := mask & p.ColorBB[p.ToMove^1]

	quiets := mask & ^p.Occupancy

	for captures != 0 {
		to := PopLSB(&captures)
		*moves = append(*moves, NewMove(sq, to, CAPTURE))
	}
	for quiets != 0 {
		to := PopLSB(&quiets)
		*moves = append(*moves, NewMove(sq, to, QUIET))
	}
}
