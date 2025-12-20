package engine

import (
	"math/bits"
)

// generates strictly legal captures and promotion moves
func (p *Position) GenTactics(moves []Move) int {
	us := p.Stm
	enemy := us ^ 1
	ksq := p.Kings[us]
	them := p.ColorOcc[enemy]

	n := 0

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

	n = p.genPawnMoves2(moves, n)

	bb := p.PieceBB[us][KNIGHT]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves2(sq, knight[sq]&them, moves, n)
	}
	bb = p.PieceBB[us][BISHOP]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves2(sq, p.pseudoBishop(sq)&them, moves, n)
	}
	bb = p.PieceBB[us][ROOK]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves2(sq, p.pseudoRook(sq)&them, moves, n)
	}
	bb = p.PieceBB[us][QUEEN]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves2(sq, (p.pseudoBishop(sq)|p.pseudoRook(sq))&them, moves, n)
	}
	return p.genKingMoves2(ksq, king[ksq]&them, moves, n)
}

func (p *Position) genKingMoves2(sq int, captures uint64, moves []Move, n int) int {
	us := p.Stm
	enemy := us ^ 1

	for captures != 0 {
		to := PopLSB(&captures)
		if p.isAttackedOcc(to, enemy, p.Occ&^(1<<sq)) {
			continue
		}
		moves[n] = NewMove(sq, to, CAPTURE)
		n++
	}

	return n
}

// generates knight and slider moves becouse they have no special cases
// pawns and kings have promotions and castling so they get their own generators
func (p *Position) genGenericMoves2(sq int, captures uint64, moves []Move, n int) int {

	if has(p.kingBlockers, sq) {
		captures &= p.allowed[sq]
	}

	for captures != 0 {
		to := PopLSB(&captures)
		moves[n] = NewMove(sq, to, CAPTURE)
		n++
	}
	return n
}

func (p *Position) genPawnMoves2(moves []Move, n int) int {

	var promo, capLeft, capRight uint64

	us := p.Stm
	enemy := us ^ 1
	empty := ^p.Occ
	P := p.PieceBB[us][PAWN]
	enemyOcc := p.ColorOcc[enemy]

	epMask := uint64(1) << p.epSquare

	if us == WHITE {
		promo = (P << 8) & empty & promotionRanks
		capLeft = ((P & notA) << 7)
		capRight = ((P & notH) << 9)
	} else {
		promo = (P >> 8) & empty & promotionRanks
		capLeft = ((P & notA) >> 9)
		capRight = ((P & notH) >> 7)
	}

	epLeft := capLeft & epMask
	epRight := capRight & epMask

	capLeft &= enemyOcc
	capRight &= enemyOcc

	promoLeft := capLeft & promotionRanks
	promoRight := capRight & promotionRanks

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
		n += 2
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
		n += 2
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
		n += 2
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
			if rookAttOcc(p.Kings[us], occ2)&(p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN]) == 0 {
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

			if rookAttOcc(p.Kings[us], occ2)&(p.PieceBB[enemy][ROOK]|p.PieceBB[enemy][QUEEN]) == 0 {
				moves[n] = NewMove(from, to, EP)
				n++
			}
		}
	}
	return n
}
