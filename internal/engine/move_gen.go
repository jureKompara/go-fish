package engine

import (
	"math/bits"
)

const promotionRanks uint64 = 0xFF000000000000FF

const notA uint64 = ^uint64(0x0101010101010101)
const notH uint64 = ^uint64(0x8080808080808080)
const rank3 uint64 = 0x0000000000FF0000
const rank6 uint64 = 0x0000FF0000000000

// returns all pseudo legal moves in the position
func (p *Position) GenMoves(moves []Move) int {
	us := p.Stm
	enemy := us ^ 1
	ksq := p.Kings[us]

	checkers := p.Checkers(ksq, enemy)

	n := 0

	switch {
	case checkers == 0:
		p.checkMask = ^uint64(0) // default: no restriction
		n = p.genCastles(moves, 0)

	case checkers&(checkers-1) == 0:
		c := bits.TrailingZeros64(checkers)
		// Knight/pawn check: ONLY capture the checker
		if has(p.PieceBB[enemy][PAWN]|p.PieceBB[enemy][KNIGHT], c) {
			p.checkMask = uint64(1) << c

		} else { // If checker is a slider, you can block OR capture
			p.checkMask = line[ksq][c]
		}

	default:
		//double check we have to move the king
		return p.genKingMoves(moves, 0)
	}

	snipers := rookAttTable[ksq][0]&(p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN]) |
		bishopAttTable[ksq][0]&(p.PieceBB[enemy][BISHOP]|p.PieceBB[enemy][QUEEN])

	p.kingBlockers = 0
	for snipers != 0 {
		sq := PopLSB(&snipers)
		betweenMask := between[ksq][sq] & p.Occ

		if betweenMask != 0 && (betweenMask&(betweenMask-1)) == 0 && betweenMask&p.ColorOcc[us] != 0 {
			p.kingBlockers |= betweenMask
			p.allowed[bits.TrailingZeros64(betweenMask)] = line[ksq][sq]
		}
	}

	mask := ^p.ColorOcc[us] & p.checkMask

	bb := p.PieceBB[us][KNIGHT]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, knight[sq]&mask, moves, n)
	}
	bb = p.PieceBB[us][BISHOP]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, p.pseudoBishop(sq)&mask, moves, n)
	}
	bb = p.PieceBB[us][ROOK]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, p.pseudoRook(sq)&mask, moves, n)
	}
	bb = p.PieceBB[us][QUEEN]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, (p.pseudoBishop(sq)|p.pseudoRook(sq))&mask, moves, n)
	}

	n = p.genKingMoves(moves, n)

	return p.genPawnMoves(moves, n)
}

func (p *Position) genKingMoves(moves []Move, n int) int {
	us := p.Stm
	enemy := us ^ 1
	ksq := p.Kings[us]

	mask := king[ksq] &^ p.ColorOcc[us]

	captures := mask & p.Occ
	quiets := mask & ^captures

	occ2 := p.Occ &^ (1 << ksq)
	for captures != 0 {
		to := PopLSB(&captures)
		if p.isAttackedOcc(to, enemy, occ2) {
			continue
		}
		moves[n] = NewMove(ksq, to, CAPTURE)
		n++
	}

	for quiets != 0 {
		to := PopLSB(&quiets)
		if p.isAttackedOcc(to, enemy, occ2) {
			continue
		}
		moves[n] = NewMove(ksq, to, QUIET)
		n++
	}
	return n
}

func (p *Position) genCastles(moves []Move, n int) int {
	us := p.Stm
	enemy := us ^ 1
	homeSquare := us*56 + 4
	//kingside castle
	if p.castleRights&(0b0001<<(2*us)) != 0 &&
		p.Occ&(0b11<<(homeSquare+1)) == 0 &&
		!p.isAttacked(homeSquare+1, enemy) &&
		!p.isAttacked(homeSquare+2, enemy) {
		moves[n] = NewMove(homeSquare, homeSquare+2, KCASTLE)
		n++
	}
	//queenside castle
	if p.castleRights&(0b0010<<(2*us)) != 0 &&
		p.Occ&(0b111<<(homeSquare-3)) == 0 &&
		!p.isAttacked(homeSquare-1, enemy) &&
		!p.isAttacked(homeSquare-2, enemy) {
		moves[n] = NewMove(homeSquare, homeSquare-2, QCASTLE)
		n++
	}
	return n
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
func bishopAttOcc(sq int, occ uint64) uint64 {
	return bishopAttTable[sq][(bishopMagic[sq]*(maskB[sq]&occ))>>bishopShifts[sq]]
}

func (p *Position) pseudoRook(sq int) uint64 {
	return rookAttTable[sq][(rookMagics[sq]*(maskR[sq]&p.Occ))>>rookShifts[sq]]
}

func (p *Position) pseudoBishop(sq int) uint64 {
	return bishopAttTable[sq][(bishopMagic[sq]*(maskB[sq]&p.Occ))>>bishopShifts[sq]]
}

// generates knight and slider moves becouse they have no special cases
// pawns and kings have promotions and castling so they get their own generators
func (p *Position) genGenericMoves(sq int, mask uint64, moves []Move, n int) int {

	if has(p.kingBlockers, sq) {
		mask &= p.allowed[sq]
	}

	captures := mask & p.Occ
	quiets := mask & ^captures

	for captures != 0 {
		to := PopLSB(&captures)
		moves[n] = NewMove(sq, to, CAPTURE)
		n++
	}
	for quiets != 0 {
		to := PopLSB(&quiets)
		moves[n] = NewMove(sq, to, QUIET)
		n++
	}
	return n
}

func (p *Position) genPawnMoves(moves []Move, n int) int {

	var singles, doubles, capLeft, capRight uint64

	us := p.Stm
	enemy := us ^ 1
	empty := ^p.Occ
	P := p.PieceBB[us][PAWN]
	enemyOcc := p.ColorOcc[enemy]

	epMask := uint64(1) << p.epSquare

	if us == WHITE {
		singles = (P << 8) & empty
		doubles = ((singles & rank3) << 8) & empty
		capLeft = ((P & notA) << 7)
		capRight = ((P & notH) << 9)
	} else {
		singles = (P >> 8) & empty
		doubles = ((singles & rank6) >> 8) & empty
		capLeft = ((P & notA) >> 9)
		capRight = ((P & notH) >> 7)
	}

	epLeft := capLeft & epMask
	epRight := capRight & epMask

	capLeft &= enemyOcc
	capRight &= enemyOcc

	singles &= p.checkMask
	doubles &= p.checkMask
	capLeft &= p.checkMask
	capRight &= p.checkMask

	promo := singles & promotionRanks
	promoLeft := capLeft & promotionRanks
	promoRight := capRight & promotionRanks

	singles ^= promo
	capLeft ^= promoLeft
	capRight ^= promoRight

	blackOffset := 16 * us

	for promoLeft != 0 {
		to := PopLSB(&promoLeft)
		from := to - 7 + blackOffset

		// pin filter (allowed indexed by from)
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, PROMOQUEENX)
		moves[n+1] = NewMove(from, to, PROMOKNIGHTX)
		moves[n+2] = NewMove(from, to, PROMOROOKX)
		moves[n+3] = NewMove(from, to, PROMOBISHOPX)
		n += 4
	}

	for promoRight != 0 {
		to := PopLSB(&promoRight)
		from := to - 9 + blackOffset

		// pin filter (allowed indexed by from)
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, PROMOQUEENX)
		moves[n+1] = NewMove(from, to, PROMOKNIGHTX)
		moves[n+2] = NewMove(from, to, PROMOROOKX)
		moves[n+3] = NewMove(from, to, PROMOBISHOPX)
		n += 4
	}

	for promo != 0 {
		to := PopLSB(&promo)
		from := to - 8 + blackOffset

		// pin filter (allowed indexed by from)
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, PROMOQUEEN)
		moves[n+1] = NewMove(from, to, PROMOKNIGHT)
		moves[n+2] = NewMove(from, to, PROMOROOK)
		moves[n+3] = NewMove(from, to, PROMOBISHOP)
		n += 4
	}

	for capLeft != 0 {
		to := PopLSB(&capLeft)
		from := to - 7 + blackOffset

		// pin filter (allowed indexed by from)
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, CAPTURE)
		n++
	}

	for capRight != 0 {
		to := PopLSB(&capRight)
		from := to - 9 + blackOffset

		// pin filter
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, CAPTURE)
		n++
	}

	if epLeft != 0 {
		to := bits.TrailingZeros64(epLeft)
		from := to - 7 + blackOffset
		capsq := from - 1

		// pin filter
		if (p.kingBlockers>>from)&1 == 0 || (p.allowed[from]>>to)&1 != 0 {
			capsqMask := uint64(1) << capsq
			occ2 := p.Occ ^ (capsqMask | uint64(1)<<from | uint64(1)<<to)
			//preven check filter
			if capsqMask&p.checkMask != 0 &&
				//rook lateral check filter
				rookAttOcc(p.Kings[us], occ2)&(p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN]) == 0 {
				moves[n] = NewMove(from, to, EP)
				n++
			}
		}
	}

	if epRight != 0 {
		to := bits.TrailingZeros64(epRight)
		from := to - 9 + blackOffset
		capsq := from + 1

		// pin filter
		if (p.kingBlockers>>from)&1 == 0 || (p.allowed[from]>>to)&1 != 0 {
			capsqMask := uint64(1) << capsq
			occ2 := p.Occ ^ (capsqMask | uint64(1)<<from | uint64(1)<<to)
			//preven check filter
			if capsqMask&p.checkMask != 0 &&
				//rook lateral check filter
				rookAttOcc(p.Kings[us], occ2)&(p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN]) == 0 {
				moves[n] = NewMove(from, to, EP)
				n++
			}
		}
	}
	for singles != 0 {
		to := PopLSB(&singles)
		from := to - 8 + blackOffset
		// pin filter
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, QUIET)
		n++
	}

	for doubles != 0 {
		to := PopLSB(&doubles)
		from := to - 16 + 2*blackOffset

		// pin filter
		if (p.kingBlockers>>from)&1 != 0 && (p.allowed[from]>>to)&1 == 0 {
			continue
		}
		moves[n] = NewMove(from, to, DOUBLE)
		n++
	}
	return n
}
