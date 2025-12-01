package main

// movement offsets for sliders
var bishOff = [4]int8{7, 9, -7, -9}
var rookOff = [4]int8{1, 8, -1, -8}
var queenOff = [8]int8{7, 9, -7, -9, 1, 8, -1, -8}

// returns all pseudo legal moves in the position
func (p *Position) pseudoAll() []Move {
	moves := make([]Move, 0, 256)
	color := p.to_move

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.pieceBB[6*p.to_move+piece]
		for bb != 0 {
			sq := popLSB(&bb)
			switch piece {
			case PAWN:
				p.GenPawnMoves(sq, color, p.pseudoPawn(sq, color), &moves)
			case KNIGHT:
				p.GenGenericMoves(sq, color, KNIGHT, p.pseudoKnight(sq, color), &moves)
			case BISHOP:
				p.GenGenericMoves(sq, color, BISHOP, p.pseudoSlider(sq, color, bishOff[:]), &moves)
			case ROOK:
				p.GenGenericMoves(sq, color, ROOK, p.pseudoSlider(sq, color, rookOff[:]), &moves)
			case QUEEN:
				p.GenGenericMoves(sq, color, QUEEN, p.pseudoSlider(sq, color, queenOff[:]), &moves)
			case KING:
				p.GenKingMoves(sq, color, p.pseudoKing(sq, color), &moves)
			}
		}
	}
	return moves
}

func (p *Position) pseudoKing(sq, color uint8) uint64 {
	return king[sq] & ^p.allBB[color]
}

func (p *Position) GenKingMoves(sq, color uint8, mask uint64, moves *[]Move) {
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
		Has(p.pieceBB[ROOK+6*color], 7+homeRank) &&
		!p.IsAttacked(sq, enemy) &&
		!p.IsAttacked(sq+1, enemy) &&
		!p.IsAttacked(sq+2, enemy) &&
		!Has(p.occupant, homeRank+5) &&
		!Has(p.occupant, homeRank+6) {
		*moves = append(*moves, NewMove(sq, sq+2, KING, NOPROMO, NOCAP, KCASTLE))
	}
	//queenside castle
	if p.castle_rights&(0b0010<<(2*color)) != 0 &&
		Has(p.pieceBB[ROOK+6*color], homeRank) &&
		!p.IsAttacked(sq, enemy) &&
		!p.IsAttacked(sq-1, enemy) &&
		!p.IsAttacked(sq-2, enemy) &&
		!Has(p.occupant, homeRank+3) &&
		!Has(p.occupant, homeRank+2) &&
		!Has(p.occupant, homeRank+1) {
		*moves = append(*moves, NewMove(sq, sq-2, KING, NOPROMO, NOCAP, QCASTLE))
	}
}

func (p *Position) pseudoPawn(sq, color uint8) uint64 {

	front := int8(sq) + 8 - 16*int8(color)
	//this should never happen
	if front < 0 || front > 63 {
		panic("Pawn wanted to jump over the edge")
	}

	//if the front square isnt empty
	if Has(p.occupant, uint8(front)) {
		return pawn[color][sq] & (p.allBB[1-color] | 1<<p.ep_square)
	}

	return (pawnPush[color][sq] & ^p.occupant) | pawn[color][sq]&(p.allBB[1-color]|1<<p.ep_square)
}

// generates pawn moves
func (p *Position) GenPawnMoves(sq, color uint8, mask uint64, moves *[]Move) {
	for mask != 0 {
		to := popLSB(&mask)
		flags := uint8(0)

		var capPiece uint8 = NOCAP

		if diff := int(to) - int(sq); diff == 16 || diff == -16 {
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

		if !(to/8 == 7 || to/8 == 0) {
			*moves = append(*moves, NewMove(sq, to, PAWN, NOPROMO, capPiece, flags))
			continue
		}
		for p := KNIGHT; p < KING; p++ {
			*moves = append(*moves, NewMove(sq, to, PAWN, p, capPiece, flags))
		}
	}
}

func (p *Position) pseudoSlider(sq, color uint8, deltas []int8) uint64 {
	var out uint64
	var prevF int8
	var sq2 int8
	for _, d := range deltas {
		sq2 = int8(sq)
		prevF = int8(sq) & 7 //this is esentialy sq%8
		for {
			sq2 += d
			if sq2 > 63 || sq2 < 0 {
				break
			}
			newF := sq2 & 7
			df := newF - prevF
			if df > 1 || df < -1 || Has(p.allBB[color], uint8(sq2)) {
				break
			}

			out |= 1 << sq2
			prevF = newF

			if Has(p.allBB[1-color], uint8(sq2)) {
				break
			}
		}
	}
	return out
}

func (p *Position) pseudoKnight(sq, color uint8) uint64 {
	return knight[sq] & ^p.allBB[color]
}

// generates knight and slide moves becouse they have no special cases
func (p *Position) GenGenericMoves(sq, color, piece uint8, mask uint64, moves *[]Move) {
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
