package engine

import (
	"strconv"
	"strings"
)

const STARTPOS string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type TestCase struct {
	FEN    string
	result []uint64
}

// fens for perft from the chess programming wiki
var Tests = []TestCase{
	{FEN: STARTPOS,
		result: []uint64{1, 20, 400, 8_902, 197_281, 4_865_609, 119_060_324, 3_195_901_860, 84_998_978_956},
	},
	{FEN: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		result: []uint64{1, 48, 2_039, 97_862, 4_085_603, 193_690_690, 8_031_647_685},
	},
	{FEN: "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		result: []uint64{1, 14, 191, 2_812, 43_238, 674_624, 11_030_083, 178_633_661, 3_009_794_393},
	},
	{FEN: "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		result: []uint64{1, 6, 264, 9_467, 422_333, 15_833_292, 706_045_033},
	},
	{FEN: "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		result: []uint64{1, 44, 1_486, 62_379, 2_103_487, 89_941_194},
	},
	{FEN: "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		result: []uint64{1, 46, 2_079, 89_890, 3_894_594, 164_075_551, 6_923_051_137, 287_188_994_746, 11_923_589_843_526},
	},
	{FEN: "r1b1k2r/ppp2qp1/2n1p2p/3pNpR1/3P4/2P3P1/PPQNP1P1/3K1B2 b kq - 1 14",
		result: []uint64{1, 31, 1073, 31872, 1077293, 31609174, 1059148863},
	},
}

// used to convert a FEN to a position
var CharToPiece = ['r' + 1]int{
	'P': PAWN,
	'N': KNIGHT,
	'B': BISHOP,
	'R': ROOK,
	'Q': QUEEN,
	'K': KING,
	'p': PAWN,
	'n': KNIGHT,
	'b': BISHOP,
	'r': ROOK,
	'q': QUEEN,
	'k': KING,
}

// used to turn a position into a FEN
var PieceToChar = [7]uint8{
	PAWN:   'P',
	KNIGHT: 'N',
	BISHOP: 'B',
	ROOK:   'R',
	QUEEN:  'Q',
	KING:   'K',
	EMPTY:  ' ',
}

func StartPos() Position {
	return FromFen(STARTPOS)
}

func FromFen(fen string) Position {
	split := strings.Split(fen, " ")
	board := split[0]
	clr := split[1]
	cr := split[2]
	ep := split[3]
	fm := split[4]
	hm := split[5]

	var pieceBB [2][6]uint64
	var colorOcc [2]uint64
	var occ uint64
	var mailbox [64]uint8
	var toMove int
	var castleRights uint8
	var epSquare int = 64 //sentinel value
	var kings [2]int

	rank := 7
	file := 0
	//board
	for _, c := range board {

		if c == '/' {
			file = 0
			rank--
			continue
		}
		if '1' <= c && c <= '8' {
			for i := range int(c - '0') {
				mailbox[rank*8+file+i] = EMPTY
			}
			file += int(c) - '0'
			continue
		}

		color := WHITE
		if 'a' <= c && c <= 'z' {
			color = BLACK
		}

		pieceType := CharToPiece[c]
		sq := rank*8 + file

		mailbox[sq] = uint8(pieceType)

		pieceBB[color][pieceType] |= (1 << sq)
		if pieceType == KING {
			kings[color] = sq
		}
		file++
	}

	//to move
	if clr == "b" {
		toMove = BLACK
	}

	//castle rights
	for _, r := range cr {
		switch r {
		case 'K':
			castleRights |= 0b0001
		case 'Q':
			castleRights |= 0b0010
		case 'k':
			castleRights |= 0b0100
		case 'q':
			castleRights |= 0b1000
		}
	}
	//ep square
	if ep != "-" {
		epSquare = int(8*(ep[1]-'1') + ep[0] - 'a')
	}

	full_move, _ := strconv.ParseInt(fm, 10, 16)

	half_move, _ := strconv.ParseInt(hm, 10, 8)

	//derived bit boards
	for piece := 0; piece <= KING; piece++ {
		colorOcc[WHITE] |= pieceBB[WHITE][piece]
		colorOcc[BLACK] |= pieceBB[BLACK][piece]
	}
	occ = colorOcc[WHITE] | colorOcc[BLACK]

	var pos = Position{
		PieceBB:      pieceBB,
		ColorOcc:     colorOcc,
		Occ:          occ,
		Board:        mailbox,
		Stm:          toMove,
		castleRights: castleRights,
		epSquare:     epSquare,
		Kings:        kings,
		fullMove:     int(full_move),
		HalfMove:     int(half_move),
	}
	pos.GenerateZobrist()
	return pos
}

func (p *Position) ExportFen() string {
	var sb strings.Builder
	var count int

	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			sq := rank*8 + file

			c := p.Board[sq]
			if c == EMPTY {
				count++
				continue
			}
			if count > 0 {
				sb.WriteByte(byte(count + '0'))
				count = 0
			}
			black := uint8(0)
			if p.ColorOcc[BLACK]&(1<<sq) != 0 {
				black = 32
			}
			sb.WriteByte(PieceToChar[c] + black)
		}
		if count > 0 {
			sb.WriteByte(byte(count + '0'))
			count = 0
		}
		if rank != 0 {
			sb.WriteByte('/')
		}
	}

	if p.Stm == 0 {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	if p.castleRights == 0 {
		sb.WriteByte('-')
	} else {
		if p.castleRights&0b0001 != 0 {
			sb.WriteByte('K')
		}
		if p.castleRights&0b0010 != 0 {
			sb.WriteByte('Q')
		}
		if p.castleRights&0b0100 != 0 {
			sb.WriteByte('k')
		}
		if p.castleRights&0b1000 != 0 {
			sb.WriteByte('q')
		}
	}
	if p.epSquare == 64 {
		sb.WriteString(" - ")
	} else {
		sb.WriteByte(' ')
		sb.WriteByte(byte('a' + p.epSquare&7))
		sb.WriteByte(byte('1' + p.epSquare>>3))
		sb.WriteByte(' ')
	}
	var buf [8]byte
	b := strconv.AppendInt(buf[:0], int64(p.fullMove), 10)
	b = append(b, ' ')
	b = strconv.AppendInt(b, int64(p.HalfMove), 10)
	sb.Write(b)

	return sb.String()
}
