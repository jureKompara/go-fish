package engine

import (
	"strconv"
	"strings"
)

const starting_pos string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type TestCase struct {
	FEN    string
	result []uint64
}

var Tests = []TestCase{
	{FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		result: []uint64{1, 20, 42069, 8902, 197281, 4865609, 119060324, 3195901860, 84_998_978_956},
	},
	{FEN: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		result: []uint64{1, 48, 2039, 97862, 4085603, 193690690, 8031647685},
	},
	{FEN: "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		result: []uint64{1, 14, 191, 2812, 43238, 674624, 11030083, 178633661, 3009794393},
	},
	{FEN: "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		result: []uint64{1, 6, 264, 9467, 422333, 15833292, 706045033},
	},
	{FEN: "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		result: []uint64{1, 44, 1486, 62379, 2103487, 89941194},
	},
	{FEN: "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		result: []uint64{1, 46, 2079, 89890, 3894594, 164075551, 6923051137, 287_188_994_746, 11_923_589_843_526},
	},
}

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

var _pieceToChar = [7]int{
	PAWN:   'P',
	KNIGHT: 'N',
	BISHOP: 'B',
	ROOK:   'R',
	QUEEN:  'Q',
	KING:   'K',
	EMPTY:  0,
}

func StartPos() Position {
	return FromFen(starting_pos)
}

func FromFen(fen string) Position {
	split := strings.Split(fen, " ")
	board := split[0]
	clr := split[1]
	cr := split[2]
	ep := split[3]
	fm := split[4]
	hm := split[5]

	var color int

	var pieceBB [12]uint64
	var allBB [2]uint64
	var occupant uint64
	var b [64]uint8
	var toMove int
	var castleRights uint8
	var epSquare int = 64 //sentinel value
	var kings [2]int
	var moveStack [512]Move
	var stateStack [512]State
	var ply int

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
				b[rank*8+file+i] = uint8(EMPTY)
			}
			file += int(c) - '0'
			continue
		}

		if 'a' <= c && c <= 'z' {
			color = 1
		} else {
			color = 0
		}

		pieceType := CharToPiece[c]
		square := rank*8 + file

		b[square] = uint8(pieceType)

		pieceBB[6*color+pieceType] |= (1 << square)
		if pieceType == KING {
			kings[color] = square
		}
		file++
	}

	//to move
	if clr == "b" {
		toMove = 1
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
		epSquare = int(8*(split[3][1]-'1') + split[3][0] - 'a')
	}

	full_move, _ := strconv.ParseInt(fm, 10, 16)

	half_move, _ := strconv.ParseInt(hm, 10, 8)

	//derived bit boards
	for i := range 6 {
		allBB[0] |= pieceBB[i]
		allBB[1] |= pieceBB[i+6]
	}
	occupant = allBB[0] | allBB[1]

	return Position{
		PieceBB:      pieceBB,
		allBB:        allBB,
		occupant:     occupant,
		Board:        b,
		ToMove:       toMove,
		castleRights: castleRights,
		epSquare:     epSquare,
		fullMove:     int(full_move),
		halfMove:     int(half_move),
		kings:        kings,
		moveStack:    moveStack,
		stateStack:   stateStack,
		Ply:          ply,
	}
}

func (p *Position) exportFen() string {
	var sb strings.Builder
	var count int

	//we build up a board to easily turn it to a fen
	var board [64]byte
	for p, bb := range p.PieceBB {
		for bb != 0 {
			board[PopLSB(&bb)] = byte(p + 1)
		}
	}

	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			sq := rank*8 + file

			c := board[sq]
			if c == 0 {
				count++
				continue
			}
			if count > 0 {
				sb.WriteByte(byte(count + '0'))
				count = 0
			}
			sb.WriteByte(byte(_pieceToChar[(c-1)%6]) + (c-1)/6*32)
		}
		if count > 0 {
			sb.WriteByte(byte(count + '0'))
			count = 0
		}
		if rank != 0 {
			sb.WriteByte('/')
		}
	}

	if p.ToMove == 0 {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	if p.castleRights == 0 {
		sb.WriteString("-")
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
	b = strconv.AppendInt(b, int64(p.halfMove), 10)
	sb.Write(b)

	return sb.String()
}
