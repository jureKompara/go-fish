package engine

import "math/bits"

// returns all pseudo legal moves in the position
func (p *Position) GenEvasions(checkers uint64) []Move {
	us := p.Stm
	enemy := us ^ 1
	ksq := p.Kings[us]

	moves := p.Movebuff[p.Ply][:]

	if checkers&(checkers-1) == 0 {
		c := bits.TrailingZeros64(checkers)
		// Knight/pawn check: ONLY capture the checker
		if has(p.PieceBB[enemy][PAWN]|p.PieceBB[enemy][KNIGHT], c) {
			p.checkMask = uint64(1) << c

		} else { // If checker is a slider, you can block OR capture
			p.checkMask = line[ksq][c]
		}

	} else {
		//double check we have to move the king
		n := p.genKingMoves(moves, 0)
		return moves[:n]
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

	n := p.genKingMoves(moves, 0)
	bb := p.PieceBB[us][KNIGHT]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, knight[sq]&p.checkMask, moves, n)
	}
	bb = p.PieceBB[us][BISHOP]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, p.pseudoBishop(sq)&p.checkMask, moves, n)
	}
	bb = p.PieceBB[us][ROOK]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, p.pseudoRook(sq)&p.checkMask, moves, n)
	}
	bb = p.PieceBB[us][QUEEN]
	for bb != 0 {
		sq := PopLSB(&bb)
		n = p.genGenericMoves(sq, (p.pseudoBishop(sq)|p.pseudoRook(sq))&p.checkMask, moves, n)
	}

	n = p.genPawnMoves3(moves, n)
	return moves[:n]
}

func (p *Position) genPawnMoves3(moves []Move, n int) int {

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
			//preven check filter
			if capsqMask == p.checkMask {
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
			//preven check filter
			if capsqMask == p.checkMask {
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
