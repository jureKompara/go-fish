package main

import (
	"strings"
)

const starting_pos string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type TestCase struct {
	FEN    string
	result []uint64
}

// fens for perft from the chess programming wiki
var Tests = []TestCase{
	{FEN: starting_pos,
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

var CharToPiece = ['r' + 1]uint8{
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

var _pieceToChar = [7]uint8{
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

	var color int

	var pieceBB [2][6]uint64
	var allBB [2]uint64
	var occupant uint64
	var b [64]uint8
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
				b[rank*8+file+i] = EMPTY
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

		pieceBB[color][pieceType] |= (1 << square)
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

	//derived bit boards
	for piece := PAWN; piece <= KING; piece++ {
		allBB[WHITE] |= pieceBB[WHITE][piece]
		allBB[BLACK] |= pieceBB[BLACK][piece]
	}
	occupant = allBB[WHITE] | allBB[BLACK]

	var pos = Position{
		PieceBB:      pieceBB,
		ColorBB:      allBB,
		Occ:          occupant,
		Board:        b,
		Stm:          toMove,
		castleRights: castleRights,
		epSquare:     uint8(epSquare),
		kings:        kings,
	}
	return pos
}
