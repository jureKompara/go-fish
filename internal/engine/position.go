package engine

// PAWN=0...KING=5,EMPTY=6
const (
	PAWN = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	EMPTY
)

const (
	WHITE = 0
	BLACK = 1
)

type Position struct {
	PieceBB      [2][6]uint64 //per piece per color bitboards
	ColorBB      [2]uint64    //per color bitboards
	Occupancy    uint64       //is there any piece here bitboard
	Board        [64]uint8    //keeps track of what piece is on Board[sq]
	ToMove       int          //side to move 0=white 1=black
	castleRights uint8        //0b1111 4 bits denote the castling rights 0123-KQkq
	epSquare     int          //denotes en passant square
	fullMove     int          //fullmove counter
	halfMove     int          //halfmove counter
	kings        [2]int       //per color king position
	//moveStack    [512]Move      //stack of move structs
	stateStack [512]State     //stack of state structs
	Ply        int            //the current Ply so we can index into the stacks
	Movebuff   [512][256]Move //buffer for storing moves per Ply
	Key        uint64         //incremental Zobrist key
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
		uint8(p.halfMove),
	}
}

func (p *Position) isAttacked(sq int, by int) bool {
	if knight[sq]&p.PieceBB[by][KNIGHT] != 0 {
		return true
	}
	if p.pseudoBishop(sq)&(p.PieceBB[by][BISHOP]|p.PieceBB[by][QUEEN]) != 0 {
		return true
	}
	if p.pseudoRook(sq)&(p.PieceBB[by][ROOK]|p.PieceBB[by][QUEEN]) != 0 {
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

func (p *Position) InCheck(stm int) bool {
	return p.isAttacked(p.kings[stm], stm^1)
}

func (p *Position) GenerateZobrist() {
	p.Key = 0
	for color := range 2 {
		for piece := PAWN; piece <= KING; piece++ {
			bb := p.PieceBB[color][piece]
			for bb != 0 {
				sq := PopLSB(&bb)
				p.Key ^= zobristPiece[color][piece][sq]
			}
		}
	}

	p.Key ^= zobristSide * uint64(p.ToMove)

	p.Key ^= zobristCastle[p.castleRights]

	if p.epSquare != 64 {
		p.Key ^= zobristEP[p.epSquare&7]
	}
}

func (p *Position) VerifyZobrist() {
	old := p.Key
	p.GenerateZobrist()
	if p.Key != old {
		panic("ZOBRIST DESYNC")
	}
}

func (p *Position) PieceAt(sq int) uint8 {
	return p.Board[sq]
}
