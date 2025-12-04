package engine

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
	pieceBB      [12]uint64     //per piece per color bitboards
	allBB        [2]uint64      //per color bitboards
	occupant     uint64         //is there any piece here bitboard
	Board        [64]uint8      //keeps track of what piece is where
	toMove       int            //side to move 0=white 1=black
	castleRights uint8          //0b1111 4 bits denote the castling rights 0123-KQkq
	epSquare     int            //denotes en passant square
	fullMove     int            //fullmove counter
	halfMove     int            //halfmove counter
	kings        [2]int         //per color king position
	moveStack    [512]Move      //stack of move structs
	stateStack   [512]State     //stack of state structs
	ply          int            //the current ply so we can index into the stacks
	movebuff     [512][256]Move //buffer for storing moves per ply
}

func (p *Position) save() {
	p.stateStack[p.ply] = State{p.castleRights, uint8(p.epSquare), uint8(p.halfMove)}
}

func (p *Position) isAttacked(sq int, by int) bool {
	if p.magicBishop(sq)&(p.pieceBB[BISHOP+6*by]|p.pieceBB[QUEEN+6*by]) != 0 {
		return true
	}
	if knight[sq]&p.pieceBB[KNIGHT+by*6] != 0 {
		return true
	}
	if p.magicRook(sq)&(p.pieceBB[ROOK+6*by]|p.pieceBB[QUEEN+6*by]) != 0 {
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

func (p *Position) whatPieceAt(sq int) int {
	return int(p.Board[sq])
}

func (p *Position) Make(move Move) {
	p.save()
	to := move.To()
	fr := move.From()
	piece := move.Piece()
	promo := move.Promo()
	flags := move.Flags()
	enemy := 1 - p.toMove

	//our piece will always end up at to
	set(&(p.pieceBB[piece+p.toMove*6]), to)
	set(&(p.allBB[p.toMove]), to)
	set(&(p.occupant), to)
	p.Board[to] = uint8(piece)

	//our piece will always leave from
	clear(&(p.pieceBB[piece+p.toMove*6]), fr)
	clear(&(p.allBB[p.toMove]), fr)
	clear(&(p.occupant), fr)
	p.Board[fr] = uint8(NOCAP)

	//if its a capture we remove enemy piece from to
	if flags&ISCAP != 0 {
		clear(&(p.pieceBB[move.Capture()+enemy*6]), to)
		clear(&(p.allBB[enemy]), to)
	}

	if p.toMove == 1 {
		p.fullMove++
	}
	p.epSquare = 64

	if piece == PAWN {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting epSquare
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8*(1-2*p.toMove)
		if flags&DP != 0 {
			p.epSquare = ep
		}
		if flags&EP != 0 {
			clear(&(p.pieceBB[PAWN+enemy*6]), ep)
			clear(&(p.allBB[enemy]), ep)
			clear(&(p.occupant), ep)
			p.Board[ep] = uint8(NOCAP)
		}
		if promo != NOPROMO {
			clear(&(p.pieceBB[PAWN+p.toMove*6]), to)
			set(&(p.pieceBB[promo+p.toMove*6]), to)
			p.Board[to] = uint8(promo)
		}

		p.halfMove = 0
	} else if flags&ISCAP != 0 {
		p.halfMove = 0
	} else {
		p.halfMove++
	}

	//castling rights when rooks move or get capped
	if to == 0 || fr == 0 {
		p.castleRights &= 0b1101
	}
	if to == 7 || fr == 7 {
		p.castleRights &= 0b1110
	}
	if to == 56 || fr == 56 {
		p.castleRights &= 0b0111
	}
	if to == 63 || fr == 63 {
		p.castleRights &= 0b1011
	}

	if piece == KING {
		p.castleRights &= 0b1100 >> (2 * p.toMove)
		p.kings[p.toMove] = to
		homeRank := p.toMove * 56
		if flags&KCASTLE != 0 {
			clear(&(p.pieceBB[ROOK+p.toMove*6]), 7+homeRank)
			clear(&(p.allBB[p.toMove]), 7+homeRank)
			clear(&(p.occupant), 7+homeRank)
			p.Board[7+homeRank] = uint8(NOCAP)

			set(&(p.pieceBB[ROOK+p.toMove*6]), 5+homeRank)
			set(&(p.allBB[p.toMove]), 5+homeRank)
			set(&(p.occupant), 5+homeRank)
			p.Board[5+homeRank] = uint8(ROOK)

		} else if flags&QCASTLE != 0 {
			clear(&(p.pieceBB[ROOK+p.toMove*6]), homeRank)
			clear(&(p.allBB[p.toMove]), homeRank)
			clear(&(p.occupant), homeRank)
			p.Board[homeRank] = uint8(NOCAP)

			set(&(p.pieceBB[ROOK+p.toMove*6]), 3+homeRank)
			set(&(p.allBB[p.toMove]), 3+homeRank)
			set(&(p.occupant), 3+homeRank)
			p.Board[3+homeRank] = uint8(ROOK)
		}
	}
	p.toMove = enemy
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
	prev_color := 1 - p.toMove

	//we always put the piece  on the from square
	set(&(p.pieceBB[piece+prev_color*6]), from)
	set(&(p.allBB[prev_color]), from)
	set(&(p.occupant), from)
	p.Board[from] = uint8(piece)

	//we clear the to square
	clear(&(p.pieceBB[piece+prev_color*6]), to)
	clear(&(p.allBB[prev_color]), to)
	clear(&(p.occupant), to)
	p.Board[to] = uint8(NOCAP)

	if promo != NOPROMO {
		clear(&(p.pieceBB[promo+prev_color*6]), to)
		set(&(p.pieceBB[PAWN+prev_color*6]), from)
	}

	//if its a cap we put the enemy piece back
	if flags&ISCAP != 0 {
		set(&(p.pieceBB[capture+p.toMove*6]), to)
		set(&(p.allBB[p.toMove]), to)
		set(&(p.occupant), to)
		p.Board[to] = uint8(capture)
	} else if flags&EP != 0 {
		behind := to + 8*(1-2*p.toMove)
		set(&(p.pieceBB[PAWN+p.toMove*6]), behind)
		set(&(p.allBB[p.toMove]), behind)
		set(&(p.occupant), behind)
		p.Board[behind] = uint8(PAWN)
	}

	if piece == KING {
		p.kings[prev_color] = from
		homeRank := prev_color * 56
		if flags&KCASTLE != 0 {
			set(&(p.pieceBB[ROOK+prev_color*6]), 7+homeRank)
			set(&(p.allBB[prev_color]), 7+homeRank)
			set(&(p.occupant), 7+homeRank)
			p.Board[7+homeRank] = uint8(ROOK)

			clear(&(p.pieceBB[ROOK+prev_color*6]), 5+homeRank)
			clear(&(p.allBB[prev_color]), 5+homeRank)
			clear(&(p.occupant), 5+homeRank)
			p.Board[5+homeRank] = uint8(NOCAP)

		} else if flags&QCASTLE != 0 {
			set(&(p.pieceBB[ROOK+prev_color*6]), homeRank)
			set(&(p.allBB[prev_color]), homeRank)
			set(&(p.occupant), homeRank)
			p.Board[homeRank] = uint8(ROOK)

			clear(&(p.pieceBB[ROOK+prev_color*6]), 3+homeRank)
			clear(&(p.allBB[prev_color]), 3+homeRank)
			clear(&(p.occupant), 3+homeRank)
			p.Board[3+homeRank] = uint8(NOCAP)
		}
	}

	//obvious stuff we can just set
	p.castleRights = state.castleRights
	p.epSquare = int(state.epSquare)
	p.halfMove = int(state.halfmove)

	if p.toMove == 0 {
		p.fullMove--
	}
	p.toMove = prev_color
	p.ply--
}
