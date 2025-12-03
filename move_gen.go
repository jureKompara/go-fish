package main

// returns all pseudo legal moves in the position
func (p *Position) pseudoAll() []Move {
	moves := p.movebuff[p.ply][:0]
	color := p.to_move

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.pieceBB[6*p.to_move+piece]
		for bb != 0 {
			sq := popLSB(&bb)
			switch piece {
			case PAWN:
				p.GenPawnMoves(sq, color, p.pseudoPawn(sq, color), &moves)
			case KNIGHT:
				p.GenGenericMoves(sq, color, KNIGHT, knight[sq] & ^p.allBB[color], &moves)
			case BISHOP:
				p.GenGenericMoves(sq, color, BISHOP, p.MagicBishop(sq) & ^p.allBB[color], &moves)
			case ROOK:
				p.GenGenericMoves(sq, color, ROOK, p.MagicRook(sq) & ^p.allBB[color], &moves)
			case QUEEN:
				p.GenGenericMoves(sq, color, QUEEN, (p.MagicRook(sq)|p.MagicBishop(sq)) & ^p.allBB[color], &moves)
			case KING:
				p.GenKingMoves(sq, color, king[sq] & ^p.allBB[color], &moves)
			}
		}
	}
	return moves
}

func (p *Position) GenKingMoves(sq int, color int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := popLSB(&mask)
		flags := uint8(0)
		capPiece := p.WhatPieceAt(to, 1-color)
		if capPiece != NOCAP {
			flags |= ISCAP
		}
		*moves = append(*moves, NewMove(sq, to, KING, NOPROMO, capPiece, flags))
	}

	enemy := 1 - color
	homeRank := color * 56
	//queenside castle
	if p.castle_rights&(0b0001<<(2*color)) != 0 &&
		!Has(p.occupant, homeRank+5) &&
		!Has(p.occupant, homeRank+6) &&
		!p.IsAttacked(sq, enemy) &&
		!p.IsAttacked(sq+1, enemy) &&
		!p.IsAttacked(sq+2, enemy) {
		*moves = append(*moves, NewMove(sq, sq+2, KING, NOPROMO, NOCAP, KCASTLE))
	}
	//queenside castle
	if p.castle_rights&(0b0010<<(2*color)) != 0 &&
		!Has(p.occupant, homeRank+3) &&
		!Has(p.occupant, homeRank+2) &&
		!Has(p.occupant, homeRank+1) &&
		!p.IsAttacked(sq, enemy) &&
		!p.IsAttacked(sq-1, enemy) &&
		!p.IsAttacked(sq-2, enemy) {
		*moves = append(*moves, NewMove(sq, sq-2, KING, NOPROMO, NOCAP, QCASTLE))
	}
}

func (p *Position) pseudoPawn(sq, color int) uint64 {
	front := sq + 8 - 16*color

	//if the front square isn't empty
	if Has(p.occupant, front) {
		return pawn[color][sq] & (p.allBB[1-color] | 1<<p.ep_square)
	}

	return (pawnPush[color][sq] & ^p.occupant) | pawn[color][sq]&(p.allBB[1-color]|1<<p.ep_square)
}

func (p *Position) GenPawnMoves(sq int, color int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := popLSB(&mask)
		flags := uint8(0)

		capPiece := NOCAP

		if diff := to - sq; diff == 16 || diff == -16 {
			flags |= DP
		}

		if to == p.ep_square {
			flags |= EP
			capPiece = PAWN
		} else { //if its not an ep we check if its a capture
			capPiece = p.WhatPieceAt(to, 1-color)
			if capPiece != NOCAP {
				flags |= ISCAP
			}
		}

		if !(to>>3 == 7 || to>>3 == 0) {
			*moves = append(*moves, NewMove(sq, to, PAWN, NOPROMO, capPiece, flags))
			continue
		}
		for p := KNIGHT; p <= QUEEN; p++ {
			*moves = append(*moves, NewMove(sq, to, PAWN, p, capPiece, flags))
		}
	}
}

func (p *Position) MagicRook(sq int) uint64 {
	index := (rookMagics[sq] * (maskR[sq] & p.occupant) >> rookShifts[sq])
	return rookAttTable[sq][index]
}

func (p *Position) MagicBishop(sq int) uint64 {
	index := (bishopMagic[sq] * (maskB[sq] & p.occupant) >> bishopShifts[sq])
	return bishopAttTable[sq][index]
}

// generates knight and slider moves becouse they have no special cases
// pawns and kings have promotions and castling so they get their own generators
func (p *Position) GenGenericMoves(sq int, color, piece int, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := popLSB(&mask)
		flags := uint8(0)
		capPiece := p.WhatPieceAt(to, 1-color)
		if capPiece != NOCAP {
			flags |= ISCAP
		}
		*moves = append(*moves, NewMove(sq, to, piece, NOPROMO, capPiece, flags))
	}
}
