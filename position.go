package main

// PAWN=0...KING=5
const (
	PAWN int = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

type Position struct {
	pieceBB       [12]uint64     //per piece per color bitboards
	allBB         [2]uint64      //per color bitboards
	occupant      uint64         //is there any piece here bitboard
	to_move       int            //side to move 0=white 1=black
	castle_rights uint8          //0b1111 4 bits denote the castling rights 0123-KQkq
	ep_square     int            //denotes en passant square
	full_move     int            //fullmove counter
	half_move     int            //halfmove counter
	kings         [2]int         //per color king position
	moveStack     [512]Move      //stack of move structs
	stateStack    [512]State     //stack of state structs
	ply           int            //the current ply so we can index into the stacks
	movebuff      [512][256]Move //buffer for storing moves per ply

}

func (p *Position) Save() {
	p.stateStack[p.ply] = State{p.castle_rights, uint8(p.ep_square), uint8(p.half_move)}
}

func (p *Position) IsAttacked(sq int, by int) bool {
	if p.MagicBishop(sq)&(p.pieceBB[BISHOP+6*by]|p.pieceBB[QUEEN+6*by]) != 0 {
		return true
	}
	if knight[sq]&p.pieceBB[KNIGHT+by*6] != 0 {
		return true
	}
	if p.MagicRook(sq)&(p.pieceBB[ROOK+6*by]|p.pieceBB[QUEEN+6*by]) != 0 {
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

func (p *Position) WhatPieceAt(sq int, color int) int {
	if !Has(p.occupant, sq) {
		return NOCAP
	}
	for piece := PAWN; piece <= KING; piece++ {
		if Has(p.pieceBB[color*6+piece], sq) {
			return piece
		}
	}
	return NOCAP
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
		clear(&(p.pieceBB[move.Capture()+enemy*6]), to)
		clear(&(p.allBB[enemy]), to)
	}

	if p.to_move == 1 {
		p.full_move++
	}
	p.ep_square = 64

	if piece == PAWN {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting ep_square
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8*(1-2*p.to_move)
		if flags&DP != 0 {
			p.ep_square = ep
		}
		if flags&EP != 0 {
			clear(&(p.pieceBB[PAWN+enemy*6]), ep)
			clear(&(p.allBB[enemy]), ep)
			clear(&(p.occupant), ep)
		}
		if promo != NOPROMO {
			clear(&(p.pieceBB[PAWN+p.to_move*6]), to)
			set(&(p.pieceBB[promo+p.to_move*6]), to)
		}

		p.half_move = 0
	} else if flags&ISCAP != 0 {
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
		homeRank := p.to_move * 56
		if flags&KCASTLE != 0 {
			clear(&(p.pieceBB[ROOK+p.to_move*6]), 7+homeRank)
			clear(&(p.allBB[p.to_move]), 7+homeRank)
			clear(&(p.occupant), 7+homeRank)

			set(&(p.pieceBB[ROOK+p.to_move*6]), 5+homeRank)
			set(&(p.allBB[p.to_move]), 5+homeRank)
			set(&(p.occupant), 5+homeRank)

		} else if flags&QCASTLE != 0 {
			clear(&(p.pieceBB[ROOK+p.to_move*6]), homeRank)
			clear(&(p.allBB[p.to_move]), homeRank)
			clear(&(p.occupant), homeRank)

			set(&(p.pieceBB[ROOK+p.to_move*6]), 3+homeRank)
			set(&(p.allBB[p.to_move]), 3+homeRank)
			set(&(p.occupant), 3+homeRank)
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
		behind := to + 8*(1-2*p.to_move)
		set(&(p.pieceBB[PAWN+p.to_move*6]), behind)
		set(&(p.allBB[p.to_move]), behind)
		set(&(p.occupant), behind)
	}

	if piece == KING {
		p.kings[prev_color] = from
		homeRank := prev_color * 56
		if flags&KCASTLE != 0 {
			set(&(p.pieceBB[ROOK+prev_color*6]), 7+homeRank)
			set(&(p.allBB[prev_color]), 7+homeRank)
			set(&(p.occupant), 7+homeRank)

			clear(&(p.pieceBB[ROOK+prev_color*6]), 5+homeRank)
			clear(&(p.allBB[prev_color]), 5+homeRank)
			clear(&(p.occupant), 5+homeRank)

		} else if flags&QCASTLE != 0 {
			set(&(p.pieceBB[ROOK+prev_color*6]), homeRank)
			set(&(p.allBB[prev_color]), homeRank)
			set(&(p.occupant), homeRank)

			clear(&(p.pieceBB[ROOK+prev_color*6]), 3+homeRank)
			clear(&(p.allBB[prev_color]), 3+homeRank)
			clear(&(p.occupant), 3+homeRank)
		}
	}

	//obvious stuff we can just set
	p.castle_rights = state.castleRights
	p.ep_square = int(state.epSquare)
	p.half_move = int(state.halfmove)

	if p.to_move == 0 {
		p.full_move--
	}
	p.to_move = prev_color
	p.ply--
}
