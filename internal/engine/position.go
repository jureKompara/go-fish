package engine

// PAWN=0...KING=5,EMPTY=6
const (
	KNIGHT = iota
	BISHOP
	ROOK
	QUEEN
	PAWN
	KING
	EMPTY
)

const (
	WHITE = 0
	BLACK = 1
)

type Position struct {
	PieceBB      [2][6]uint64 //per piece per color bitboards
	ColorOcc     [2]uint64    //per color bitboards
	Occ          uint64       //is there any piece here bitboard
	Board        [64]uint8    //keeps track of what piece is on Board[sq]
	Stm          int          //side to move 0=white 1=black
	castleRights uint8        //0b1111 4 bits denote the castling rights 0123-KQkq
	epSquare     int          //denotes en passant square
	Kings        [2]int       //per color king position

	Ply         int            //the current Ply so we can index into the stacks
	stateStack  [512]State     //stack of state structs
	Movebuff    [512][256]Move //buffer for storing moves per Ply
	Hash        uint64         //incremental Zobrist hash key
	HashHistory [512]uint64    //a history of hashes for 3fold rep

	fullMove int //fullmove counter
	HalfMove int //halfmove counter

	kingBlockers uint64
	allowed      [64]uint64
	checkMask    uint64
}

var zobristPiece [2][6][64]uint64
var zobristSide uint64
var zobristCastle [16]uint64
var zobristEP [8]uint64

func (p *Position) save(capture uint8) {
	p.stateStack[p.Ply] = State{
		capture,
		p.castleRights,
		uint8(p.epSquare),
		uint8(p.HalfMove),
		p.Hash,
	}
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

func (p *Position) InCheck() bool {
	return p.isAttacked(p.Kings[p.Stm], p.Stm^1)
}

func (p *Position) GenerateZobrist() {
	p.Hash = 0
	for color := range 2 {
		for piece := 0; piece <= KING; piece++ {
			bb := p.PieceBB[color][piece]
			for bb != 0 {
				sq := PopLSB(&bb)
				p.Hash ^= zobristPiece[color][piece][sq]
			}
		}
	}

	p.Hash ^= zobristSide * uint64(p.Stm)

	p.Hash ^= zobristCastle[p.castleRights]

	if p.epSquare != 64 {
		p.Hash ^= zobristEP[p.epSquare&7]
	}
}

func (p *Position) VerifyZobrist() {
	old := p.Hash
	p.GenerateZobrist()
	if p.Hash != old {
		panic("ZOBRIST DESYNC")
	}
}
