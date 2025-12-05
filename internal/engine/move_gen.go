package engine

// returns all pseudo legal moves in the position
func (p *Position) GenMoves() []Move {
	moves := p.Movebuff[p.Ply][:0]
	color := p.ToMove

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.PieceBB[6*p.ToMove+piece]
		for bb != 0 {
			sq := PopLSB(&bb)
			switch piece {
			case PAWN:
				p.genPawnMoves(sq, p.pseudoPawn(sq, color), &moves)
			case KNIGHT:
				p.genGenericMoves(sq, KNIGHT, knight[sq] & ^p.allBB[color], &moves)
			case BISHOP:
				p.genGenericMoves(sq, BISHOP, p.magicBishop(sq) & ^p.allBB[color], &moves)
			case ROOK:
				p.genGenericMoves(sq, ROOK, p.magicRook(sq) & ^p.allBB[color], &moves)
			case QUEEN:
				p.genGenericMoves(sq, QUEEN, (p.magicRook(sq)|p.magicBishop(sq)) & ^p.allBB[color], &moves)
			case KING:
				p.genKingMoves(sq, color, king[sq] & ^p.allBB[color], &moves)
			}
		}
	}
	return moves
}

func (p *Position) genKingMoves(sq int, color int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := PopLSB(&mask)
		flags := uint8(0)
		capPiece := int(p.Board[to])
		if capPiece != EMPTY {
			flags |= ISCAP
		}
		*moves = append(*moves, NewMove(sq, to, KING, EMPTY, capPiece, flags))
	}

	enemy := 1 - color
	homeRank := color * 56
	//queenside castle
	if p.castleRights&(0b0001<<(2*color)) != 0 &&
		!has(p.occupant, homeRank+5) &&
		!has(p.occupant, homeRank+6) &&
		!p.isAttacked(sq, enemy) &&
		!p.isAttacked(sq+1, enemy) &&
		!p.isAttacked(sq+2, enemy) {
		*moves = append(*moves, NewMove(sq, sq+2, KING, EMPTY, EMPTY, KCASTLE))
	}
	//queenside castle
	if p.castleRights&(0b0010<<(2*color)) != 0 &&
		!has(p.occupant, homeRank+3) &&
		!has(p.occupant, homeRank+2) &&
		!has(p.occupant, homeRank+1) &&
		!p.isAttacked(sq, enemy) &&
		!p.isAttacked(sq-1, enemy) &&
		!p.isAttacked(sq-2, enemy) {
		*moves = append(*moves, NewMove(sq, sq-2, KING, EMPTY, EMPTY, QCASTLE))
	}
}

func (p *Position) pseudoPawn(sq, color int) uint64 {
	front := sq + 8 - 16*color
	//if the front square isn't empty
	if has(p.occupant, front) {
		return pawn[color][sq] & (p.allBB[1-color] | 1<<p.epSquare)
	}
	return (pawnPush[color][sq] & ^p.occupant) | pawn[color][sq]&(p.allBB[1-color]|1<<p.epSquare)
}

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

func (p *Position) magicRook(sq int) uint64 {
	index := (rookMagics[sq] * (maskR[sq] & p.occupant) >> rookShifts[sq])
	return rookAttTable[sq][index]
}

func (p *Position) magicBishop(sq int) uint64 {
	index := (bishopMagic[sq] * (maskB[sq] & p.occupant) >> bishopShifts[sq])
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
