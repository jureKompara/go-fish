package main

// PAWN=0...KING=5,EMPTY=6
const (
	PAWN uint8 = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	EMPTY
)

const (
	WHITE uint8 = 0
	BLACK uint8 = 1
)

type Position struct {
	PieceBB      [2][6]uint64   //per piece per color bitboards
	ColorOcc     [2]uint64      //per color bitboards
	Occ          uint64         //is there any piece here bitboard
	Board        [64]uint8      //keeps track of what piece is on Board[sq]
	Stm          int            //side to move 0=white 1=black
	castleRights uint8          //0b1111 4 bits denote the castling rights 0123-KQkq
	epSquare     uint8          //denotes en passant square
	kings        [2]int         //per color king position
	stateStack   [512]State     //stack of state structs
	Ply          int            //the current Ply so we can index into the stacks
	Movebuff     [512][256]Move //buffer for storing moves per Ply
	kingBlockers uint64
	allowed      [64]uint64
	checkMask    uint64
}

func (p *Position) save(capture uint8) {
	p.stateStack[p.Ply] = State(capture|p.castleRights<<3) | State(p.epSquare)<<7

}

func (p *Position) isAttacked(sq int, by int) bool {
	if p.pseudoBishop(sq)&(p.PieceBB[by][BISHOP]|p.PieceBB[by][QUEEN]) != 0 {
		return true
	}
	if p.pseudoRook(sq)&(p.PieceBB[by][ROOK]|p.PieceBB[by][QUEEN]) != 0 {
		return true
	}
	if knight[sq]&p.PieceBB[by][KNIGHT] != 0 {
		return true
	}
	if pawn[by^1][sq]&p.PieceBB[by][PAWN] != 0 {
		return true
	}
	if king[sq]&p.PieceBB[by][KING] != 0 {
		return true
	}
	return false
}

func (p *Position) isAttackedOcc(sq int, by int, occ uint64) bool {
	if bishopAttOcc(sq, occ)&(p.PieceBB[by][BISHOP]|p.PieceBB[by][QUEEN]) != 0 {
		return true
	}
	if rookAttOcc(sq, occ)&(p.PieceBB[by][ROOK]|p.PieceBB[by][QUEEN]) != 0 {
		return true
	}
	if knight[sq]&p.PieceBB[by][KNIGHT] != 0 {
		return true
	}
	if pawn[by^1][sq]&p.PieceBB[by][PAWN] != 0 {
		return true
	}
	if king[sq]&p.PieceBB[by][KING] != 0 {
		return true
	}
	return false
}
