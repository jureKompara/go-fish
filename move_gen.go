package main

import (
	"math/bits"
)

// used for promo only generation for generating tactical moves
const promotionRanks uint64 = 0xFF000000000000FF

// returns all pseudo legal moves in the position
func (p *Position) GenMoves(moves *[]Move) {
	us := p.Stm
	enemy := us ^ 1
	ksq := p.kings[us]

	checkers := p.Checkers(ksq, enemy)
	checkCount := bits.OnesCount64(checkers)

	if checkCount == 2 {
		p.genKingMoves(ksq, moves)
		return
	}

	p.checkMask = ^uint64(0) // default: no restriction
	if checkCount == 1 {
		c := bits.TrailingZeros64(checkers)

		// If checker is a slider, you can block OR capture
		if has(p.PieceBB[enemy][BISHOP]|p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN], c) {
			p.checkMask = line[ksq][c]
		} else {
			// Knight/pawn check: ONLY capture the checker
			p.checkMask = uint64(1) << c
		}
	}

	snipers := rookAttTable[ksq][0]&(p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN]) |
		bishopAttTable[ksq][0]&(p.PieceBB[enemy][BISHOP]|p.PieceBB[enemy][QUEEN])

	p.kingBlockers = 0
	for snipers != 0 {
		sq := PopLSB(&snipers)
		betweenMask := between[ksq][sq] & p.Occ
		if bits.OnesCount64(betweenMask) == 1 && betweenMask&p.ColorBB[us] != 0 {
			p.kingBlockers |= betweenMask
			p.allowed[bits.TrailingZeros64(betweenMask)] = line[ksq][sq]
		}
	}

	for piece := PAWN; piece <= KING; piece++ {
		bb := p.PieceBB[p.Stm][piece]
		for bb != 0 {
			sq := PopLSB(&bb)
			switch piece {
			case PAWN:
				p.genPawnMoves(sq, moves)
			case KNIGHT:
				p.genGenericMoves(sq, knight[sq] & ^p.ColorBB[us], moves)
			case BISHOP:
				p.genGenericMoves(sq, p.pseudoBishop(sq) & ^p.ColorBB[us], moves)
			case ROOK:
				p.genGenericMoves(sq, p.pseudoRook(sq) & ^p.ColorBB[us], moves)
			case QUEEN:
				p.genGenericMoves(sq, (p.pseudoBishop(sq)|p.pseudoRook(sq)) & ^p.ColorBB[us], moves)
			case KING:
				p.genKingMoves(sq, moves)
			}
		}
	}
}

func (p *Position) genKingMoves(sq int, moves *[]Move) {
	us := p.Stm
	enemy := us ^ 1

	mask := king[sq] & ^p.ColorBB[us]

	captures := mask & p.Occ
	quiets := mask & ^captures

	clear(&p.Occ, sq)

	for captures != 0 {
		to := PopLSB(&captures)
		if p.isAttacked(to, enemy) {
			continue
		}
		*moves = append(*moves, NewMove(sq, to, CAPTURE))
	}

	for quiets != 0 {
		to := PopLSB(&quiets)
		if p.isAttacked(to, enemy) {
			continue
		}
		*moves = append(*moves, NewMove(sq, to, QUIET))
	}

	set(&p.Occ, sq)

	//i already have this info...should do better
	if !p.isAttacked(sq, enemy) {
		homeRank := us * 56
		//kingside castle
		if p.castleRights&(0b0001<<(2*us)) != 0 &&
			!has(p.Occ, homeRank+5) &&
			!has(p.Occ, homeRank+6) &&
			!p.isAttacked(sq+1, enemy) &&
			!p.isAttacked(sq+2, enemy) {
			*moves = append(*moves, NewMove(sq, sq+2, KCASTLE))
		}
		//queenside castle
		if p.castleRights&(0b0010<<(2*us)) != 0 &&
			!has(p.Occ, homeRank+3) &&
			!has(p.Occ, homeRank+2) &&
			!has(p.Occ, homeRank+1) &&
			!p.isAttacked(sq-1, enemy) &&
			!p.isAttacked(sq-2, enemy) {
			*moves = append(*moves, NewMove(sq, sq-2, QCASTLE))
		}
	}
}

// returns a bitboard with all posible pawn moves
func (p *Position) pseudoPawn(sq int) uint64 {
	us := p.Stm
	enemy := us ^ 1
	front := sq + 8 - 16*us
	//if the front square is occupied
	if has(p.Occ, front) {
		return pawn[us][sq] & (p.ColorBB[enemy] | 1<<p.epSquare)
	}
	return (pawnPush[us][sq] & ^p.Occ) | pawn[us][sq]&(p.ColorBB[enemy]|1<<p.epSquare)
}

func (p *Position) genPawnMoves(sq int, moves *[]Move) {
	us := p.Stm
	enemy := us ^ 1
	mask := p.pseudoPawn(sq)

	if has(p.kingBlockers, sq) {
		mask &= p.allowed[sq]
	}

	ep := mask & (1 << p.epSquare)

	mask &= p.checkMask

	captures := mask & p.Occ

	quiets := mask & ^(captures | ep)

	promoCaptures := captures & promotionRanks
	promoQuiets := quiets & promotionRanks

	captures &= ^promoCaptures
	quiets &= ^promoQuiets

	for promoCaptures != 0 {
		to := PopLSB(&promoCaptures)
		*moves = append(*moves, NewMove(sq, to, PROMOQUEENX))
		*moves = append(*moves, NewMove(sq, to, PROMOKNIGHTX))
		*moves = append(*moves, NewMove(sq, to, PROMOROOKX))
		*moves = append(*moves, NewMove(sq, to, PROMOBISHOPX))
	}

	for promoQuiets != 0 {
		to := PopLSB(&promoQuiets)
		*moves = append(*moves, NewMove(sq, to, PROMOQUEEN))
		*moves = append(*moves, NewMove(sq, to, PROMOKNIGHT))
		*moves = append(*moves, NewMove(sq, to, PROMOROOOK))
		*moves = append(*moves, NewMove(sq, to, PROMOBISHOP))
	}

	for captures != 0 {
		to := PopLSB(&captures)
		*moves = append(*moves, NewMove(sq, to, CAPTURE))
	}

	for quiets != 0 {
		to := PopLSB(&quiets)
		if diff := to - sq; diff == 16 || diff == -16 {
			*moves = append(*moves, NewMove(sq, to, DOUBLE))
			continue
		}
		*moves = append(*moves, NewMove(sq, to, QUIET))
	}

	if ep != 0 {
		to := bits.TrailingZeros64(ep)
		capSq := to - 8*(1-2*us) // white: to-8, black: to+8

		occ2 := p.Occ ^ (uint64(1)<<sq | uint64(1)<<capSq | uint64(1)<<to)

		if (rookAttOcc(p.kings[us], occ2) & (p.PieceBB[enemy][ROOK] | p.PieceBB[enemy][QUEEN])) == 0 {
			*moves = append(*moves, NewMove(sq, to, EP))
		}
	}
}

func (p *Position) Checkers(sq int, by int) uint64 {
	return p.pseudoBishop(sq)&(p.PieceBB[by][BISHOP]|p.PieceBB[by][QUEEN]) |
		p.pseudoRook(sq)&(p.PieceBB[by][ROOK]|p.PieceBB[by][QUEEN]) |
		knight[sq]&p.PieceBB[by][KNIGHT] |
		pawn[by^1][sq]&p.PieceBB[by][PAWN]
}

func rookAttOcc(sq int, occ uint64) uint64 {
	return rookAttTable[sq][(rookMagics[sq]*(maskR[sq]&occ))>>rookShifts[sq]]
}

func (p *Position) pseudoRook(sq int) uint64 {
	return rookAttTable[sq][(rookMagics[sq] * (maskR[sq] & p.Occ) >> rookShifts[sq])]
}

func (p *Position) pseudoBishop(sq int) uint64 {
	return bishopAttTable[sq][(bishopMagic[sq] * (maskB[sq] & p.Occ) >> bishopShifts[sq])]
}

// generates knight and slider moves becouse they have no special cases
// pawns and kings have promotions and castling so they get their own generators
func (p *Position) genGenericMoves(sq int, mask uint64, moves *[]Move) {

	mask &= p.checkMask

	if has(p.kingBlockers, sq) {
		mask &= p.allowed[sq]
	}

	captures := mask & p.Occ
	quiets := mask & ^captures

	for captures != 0 {
		to := PopLSB(&captures)
		*moves = append(*moves, NewMove(sq, to, CAPTURE))
	}
	for quiets != 0 {
		to := PopLSB(&quiets)
		*moves = append(*moves, NewMove(sq, to, QUIET))
	}
}
