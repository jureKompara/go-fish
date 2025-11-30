package main

// PAWN=0...KING=5
const (
	PAWN uint8 = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

type Position struct {
	pieceBB       [12]uint64 //per piece per color bitboards
	allBB         [2]uint64  //per color bitboards
	occupant      uint64     //is there any piece here ass bitboard
	to_move       uint8      //side to move 0=white 1=black
	castle_rights uint8      //0b1111 4 bits denote the castling rights 0123-KQkq
	ep_square     uint8      //denotes en passant square
	full_move     uint16     //fullmove counter
	half_move     uint8      //halfmove counter
	kings         [2]uint8   //per color king position
	moveStack     [512]Move  //stock of move structs
	stateStack    [512]State //stack of state structs
	ply           uint16     //the current ply so we can index into the stacks
}

func (p *Position) Save() {
	p.stateStack[p.ply] = State{p.castle_rights, p.ep_square, p.half_move}
}

func (p *Position) IsAttacked(sq, by uint8) bool {
	if knight[sq]&p.pieceBB[KNIGHT+by*6] != 0 {
		return true
	}
	if p.pseudoSlider(sq, 1-by, rookOff[:])&
		(p.pieceBB[ROOK+by*6]|p.pieceBB[QUEEN+by*6]) != 0 {
		return true
	}
	if p.pseudoSlider(sq, 1-by, bishOff[:])&
		(p.pieceBB[BISHOP+by*6]|p.pieceBB[QUEEN+by*6]) != 0 {
		return true
	}
	if pawn[1-by][sq]&p.pieceBB[PAWN+by*6] != 0 {
		return true
	}
	if king[sq]&p.pieceBB[KING+by*6] != 0 {
		return true
	}
	return false
}

func (p *Position) WhatPieceAt(sq, color uint8) (uint8, bool) {
	for piece := PAWN; piece <= KING; piece++ {
		if Has(p.pieceBB[color*6+piece], sq) {
			return piece, true
		}
	}
	return NOCAP, false
}

func (p *Position) Make(move Move) {
	p.Save()
	to := move.To()
	fr := move.From()
	piece := move.Piece()
	promo := move.Promo()
	flags := move.Flags()
	enemy := 1 - p.to_move

	//our piece will always end up at to
	set(&(p.pieceBB[piece+p.to_move*6]), to)
	set(&(p.allBB[p.to_move]), to)
	set(&(p.occupant), to)

	//our piece will always leave from
	clear(&(p.pieceBB[piece+p.to_move*6]), fr)
	clear(&(p.allBB[p.to_move]), fr)
	clear(&(p.occupant), fr)

	//if its a capture we remove enemy piece from to
	if flags&ISCAP != 0 {
		clear(&(p.pieceBB[move.Capture()+(1-p.to_move)*6]), to)
		clear(&(p.allBB[1-p.to_move]), to)
	}

	if p.to_move == 1 {
		p.full_move++
	}
	p.ep_square = 64
	if piece == PAWN || flags&ISCAP != 0 {
		if flags&DP != 0 {
			p.ep_square = uint8(int(to) - 8*(1-2*int(p.to_move)))
		}
		if flags&EP != 0 {
			clear(&(p.pieceBB[PAWN+(enemy)*6]), to-8*(1-2*p.to_move))
			clear(&(p.allBB[enemy]), to-8*(1-2*p.to_move))
			clear(&(p.occupant), to-8*(1-2*p.to_move))
		}
		if promo != NOPROMO {
			clear(&(p.pieceBB[PAWN+p.to_move*6]), to)
			set(&(p.pieceBB[promo+p.to_move*6]), to)
		}

		p.half_move = 0
	} else {
		p.half_move++
	}

	//castling rights when rooks move or get capped
	if to == 0 || fr == 0 {
		p.castle_rights &= 0b1101
	}
	if to == 7 || fr == 7 {
		p.castle_rights &= 0b1110
	}
	if to == 56 || fr == 56 {
		p.castle_rights &= 0b0111
	}
	if to == 63 || fr == 63 {
		p.castle_rights &= 0b1011
	}

	if piece == KING {
		p.castle_rights &= 0b1100 >> (2 * p.to_move)
		p.kings[p.to_move] = to

		if flags&KCASTLE != 0 {
			clear(&(p.pieceBB[ROOK+p.to_move*6]), 7+p.to_move*7*8)
			clear(&(p.allBB[p.to_move]), 7+p.to_move*7*8)
			clear(&(p.occupant), 7+p.to_move*7*8)

			set(&(p.pieceBB[ROOK+p.to_move*6]), 5+p.to_move*7*8)
			set(&(p.allBB[p.to_move]), 5+p.to_move*7*8)
			set(&(p.occupant), 5+p.to_move*7*8)

		} else if flags&QCASTLE != 0 {
			clear(&(p.pieceBB[ROOK+p.to_move*6]), p.to_move*7*8)
			clear(&(p.allBB[p.to_move]), p.to_move*7*8)
			clear(&(p.occupant), p.to_move*7*8)

			set(&(p.pieceBB[ROOK+p.to_move*6]), 3+p.to_move*7*8)
			set(&(p.allBB[p.to_move]), 3+p.to_move*7*8)
			set(&(p.occupant), 3+p.to_move*7*8)
		}
	}
	p.to_move = enemy
	p.ply++
}

func (p *Position) Unmake(move Move) {
	state := p.stateStack[p.ply-1]
	from := move.From()
	to := move.To()
	piece := move.Piece()
	capture := move.Capture()
	promo := move.Promo()
	flags := move.Flags()
	prev_color := 1 - p.to_move

	//we always put the piece  on the from square
	set(&(p.pieceBB[piece+prev_color*6]), from)
	set(&(p.allBB[prev_color]), from)
	set(&(p.occupant), from)

	//we clear the to square
	clear(&(p.pieceBB[piece+prev_color*6]), to)
	clear(&(p.allBB[prev_color]), to)
	clear(&(p.occupant), to)

	if promo != NOPROMO {
		clear(&(p.pieceBB[promo+prev_color*6]), to)
		set(&(p.pieceBB[PAWN+prev_color*6]), from)
	}

	//if its a cap we put the enemy piece back
	if flags&ISCAP != 0 {
		set(&(p.pieceBB[capture+p.to_move*6]), to)
		set(&(p.allBB[p.to_move]), to)
		set(&(p.occupant), to)
	} else if flags&EP != 0 {
		set(&(p.pieceBB[PAWN+p.to_move*6]), to+8*(1-2*p.to_move))
		set(&(p.allBB[p.to_move]), to+8*(1-2*p.to_move))
		set(&(p.occupant), to+8*(1-2*p.to_move))
	}

	if piece == KING {
		p.kings[prev_color] = from
		if flags&KCASTLE != 0 {
			set(&(p.pieceBB[ROOK+prev_color*6]), 7+prev_color*7*8)
			set(&(p.allBB[prev_color]), 7+prev_color*7*8)
			set(&(p.occupant), 7+prev_color*7*8)

			clear(&(p.pieceBB[ROOK+prev_color*6]), 5+prev_color*7*8)
			clear(&(p.allBB[prev_color]), 5+prev_color*7*8)
			clear(&(p.occupant), 5+prev_color*7*8)

		} else if flags&QCASTLE != 0 {
			set(&(p.pieceBB[ROOK+prev_color*6]), prev_color*7*8)
			set(&(p.allBB[prev_color]), prev_color*7*8)
			set(&(p.occupant), prev_color*7*8)

			clear(&(p.pieceBB[ROOK+prev_color*6]), 3+prev_color*7*8)
			clear(&(p.allBB[prev_color]), 3+prev_color*7*8)
			clear(&(p.occupant), 3+prev_color*7*8)
		}
	}

	//obvious stuff we can just set
	p.castle_rights = state.castleRights
	p.ep_square = state.epSquare
	p.half_move = state.halfmove

	if p.to_move == 0 {
		p.full_move--
	}
	p.to_move = prev_color
	p.ply--
}
