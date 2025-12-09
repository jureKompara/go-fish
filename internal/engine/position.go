package engine

// PAWN=0...KING=5
const (
	PAWN int = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	EMPTY
)

const (
	WHITE = iota
	BLACK
)

type Position struct {
	PieceBB      [12]uint64     //per piece per color bitboards
	ColorBB      [2]uint64      //per color bitboards
	Occupancy    uint64         //is there any piece here bitboard
	Board        [64]uint8      //keeps track of what piece is on Board[sq]
	ToMove       int            //side to move 0=white 1=black
	castleRights uint8          //0b1111 4 bits denote the castling rights 0123-KQkq
	epSquare     int            //denotes en passant square
	fullMove     int            //fullmove counter
	halfMove     int            //halfmove counter
	kings        [2]int         //per color king position
	moveStack    [512]Move      //stack of move structs
	stateStack   [512]State     //stack of state structs
	Ply          int            //the current Ply so we can index into the stacks
	Movebuff     [512][256]Move //buffer for storing moves per Ply
}

func (p *Position) save() {
	p.stateStack[p.Ply] = State{p.castleRights, uint8(p.epSquare), uint8(p.halfMove)}
}

func (p *Position) isAttacked(sq int, by int) bool {
	if p.pseudoBishop(sq)&(p.PieceBB[BISHOP+6*by]|p.PieceBB[QUEEN+6*by]) != 0 {
		return true
	}
	if knight[sq]&p.PieceBB[KNIGHT+by*6] != 0 {
		return true
	}
	if p.pseudoRook(sq)&(p.PieceBB[ROOK+6*by]|p.PieceBB[QUEEN+6*by]) != 0 {
		return true
	}
	if pawn[by^1][sq]&p.PieceBB[PAWN+by*6] != 0 {
		return true
	}
	if king[sq]&p.PieceBB[KING+by*6] != 0 {
		return true
	}
	return false
}

func (p *Position) InCheck(stm int) bool {
	//fmt.Println("king:", 1-p.toMove, "is attacked by", p.toMove)
	return p.isAttacked(p.kings[stm], stm^1)
}

func (p *Position) Make(move Move) {
	p.save()
	to := move.To()
	fr := move.From()
	piece := move.Piece()
	promo := move.Promo()
	flags := move.Flags()
	us := p.ToMove
	enemy := 1 - us

	//our piece will always end up at to
	set(&(p.PieceBB[piece+us*6]), to)
	set(&(p.ColorBB[us]), to)
	set(&(p.Occupancy), to)
	p.Board[to] = uint8(piece)

	//our piece will always leave from
	clear(&(p.PieceBB[piece+us*6]), fr)
	clear(&(p.ColorBB[us]), fr)
	clear(&(p.Occupancy), fr)
	p.Board[fr] = uint8(EMPTY)

	//if its a capture we remove enemy piece from to
	if flags&ISCAP != 0 {
		clear(&(p.PieceBB[move.Capture()+enemy*6]), to)
		clear(&(p.ColorBB[enemy]), to)
	}

	if us == 1 {
		p.fullMove++
	}
	p.epSquare = 64

	if piece == PAWN {
		//for both double pushes and ep captures the relevant squre
		//is the one behind 'to' so we can use it for setting epSquare
		//after a double push or clearing the pawn after an ep capture
		ep := to - 8*(1-2*us)
		if flags&DP != 0 {
			p.epSquare = ep
		}
		if flags&EP != 0 {
			clear(&(p.PieceBB[PAWN+enemy*6]), ep)
			clear(&(p.ColorBB[enemy]), ep)
			clear(&(p.Occupancy), ep)
			p.Board[ep] = uint8(EMPTY)
		}
		if promo != EMPTY {
			clear(&(p.PieceBB[PAWN+us*6]), to)
			set(&(p.PieceBB[promo+us*6]), to)
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
		p.castleRights &= 0b1100 >> (2 * us)
		p.kings[us] = to
		homeRank := us * 56
		if flags&KCASTLE != 0 {
			clear(&(p.PieceBB[ROOK+us*6]), 7+homeRank)
			clear(&(p.ColorBB[us]), 7+homeRank)
			clear(&(p.Occupancy), 7+homeRank)
			p.Board[7+homeRank] = uint8(EMPTY)

			set(&(p.PieceBB[ROOK+us*6]), 5+homeRank)
			set(&(p.ColorBB[us]), 5+homeRank)
			set(&(p.Occupancy), 5+homeRank)
			p.Board[5+homeRank] = uint8(ROOK)

		} else if flags&QCASTLE != 0 {
			clear(&(p.PieceBB[ROOK+us*6]), homeRank)
			clear(&(p.ColorBB[us]), homeRank)
			clear(&(p.Occupancy), homeRank)
			p.Board[homeRank] = uint8(EMPTY)

			set(&(p.PieceBB[ROOK+us*6]), 3+homeRank)
			set(&(p.ColorBB[us]), 3+homeRank)
			set(&(p.Occupancy), 3+homeRank)
			p.Board[3+homeRank] = uint8(ROOK)
		}
	}
	p.ToMove = enemy
	p.Ply++
}

func (p *Position) Unmake(move Move) {
	state := p.stateStack[p.Ply-1]
	from := move.From()
	to := move.To()
	piece := move.Piece()
	capture := move.Capture()
	promo := move.Promo()
	flags := move.Flags()

	enemy := p.ToMove
	us := 1 - enemy

	//we always put the piece  on the from square
	set(&(p.PieceBB[piece+us*6]), from)
	set(&(p.ColorBB[us]), from)
	set(&(p.Occupancy), from)
	p.Board[from] = uint8(piece)

	//we clear the to square
	clear(&(p.PieceBB[piece+us*6]), to)
	clear(&(p.ColorBB[us]), to)
	clear(&(p.Occupancy), to)
	p.Board[to] = uint8(EMPTY)

	if promo != EMPTY {
		clear(&(p.PieceBB[promo+us*6]), to)
		set(&(p.PieceBB[PAWN+us*6]), from)
	}

	//if its a cap we put the enemy piece back
	if flags&ISCAP != 0 {
		set(&(p.PieceBB[capture+enemy*6]), to)
		set(&(p.ColorBB[enemy]), to)
		set(&(p.Occupancy), to)
		p.Board[to] = uint8(capture)
	} else if flags&EP != 0 {
		behind := to + 8*(1-2*enemy)
		set(&(p.PieceBB[PAWN+enemy*6]), behind)
		set(&(p.ColorBB[enemy]), behind)
		set(&(p.Occupancy), behind)
		p.Board[behind] = uint8(PAWN)
	}

	if piece == KING {
		p.kings[us] = from
		homeRank := us * 56
		if flags&KCASTLE != 0 {
			set(&(p.PieceBB[ROOK+us*6]), 7+homeRank)
			set(&(p.ColorBB[us]), 7+homeRank)
			set(&(p.Occupancy), 7+homeRank)
			p.Board[7+homeRank] = uint8(ROOK)

			clear(&(p.PieceBB[ROOK+us*6]), 5+homeRank)
			clear(&(p.ColorBB[us]), 5+homeRank)
			clear(&(p.Occupancy), 5+homeRank)
			p.Board[5+homeRank] = uint8(EMPTY)

		} else if flags&QCASTLE != 0 {
			set(&(p.PieceBB[ROOK+us*6]), homeRank)
			set(&(p.ColorBB[us]), homeRank)
			set(&(p.Occupancy), homeRank)
			p.Board[homeRank] = uint8(ROOK)

			clear(&(p.PieceBB[ROOK+us*6]), 3+homeRank)
			clear(&(p.ColorBB[us]), 3+homeRank)
			clear(&(p.Occupancy), 3+homeRank)
			p.Board[3+homeRank] = uint8(EMPTY)
		}
	}

	//obvious stuff we can just set
	p.castleRights = state.castleRights
	p.epSquare = int(state.epSquare)
	p.halfMove = int(state.halfmove)

	if enemy == WHITE {
		p.fullMove--
	}
	p.ToMove = us
	p.Ply--
}
