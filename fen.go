package main

import (
	"strconv"
	"strings"
)

var charToUint8 = [...]uint8{
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

var Uint8ToChar = [...]uint8{
	PAWN:   'P',
	KNIGHT: 'N',
	BISHOP: 'B',
	ROOK:   'R',
	QUEEN:  'Q',
	KING:   'K',
}

func FromFen(fen string) Position {
	var split []string = strings.Split(fen, " ")

	var rank uint8 = 7
	var file uint8 = 0

	var c uint8
	var pieceType uint8
	var square uint8
	var color uint8

	var pieceBB [12]uint64
	var allBB [2]uint64
	var occupant uint64
	var to_move uint8
	var castle_rights uint8
	var ep_square uint8 = 64 //sentinel value
	var full_move uint16
	var half_move uint8
	var kings [2]uint8
	var moveStack [512]Move
	var stateStack [512]State
	var ply uint16

	//board
	for i := range split[0] {
		c = split[0][i]

		if c == '/' {
			file = 0
			rank--
			continue
		}
		if '1' <= c && c <= '8' {
			file += c - '0'
			continue
		}

		if 'a' <= c && c <= 'z' {
			color = 1
		} else {
			color = 0
		}

		pieceType = charToUint8[c]
		square = rank*8 + file

		pieceBB[6*color+pieceType] |= (1 << square)
		if pieceType == KING {
			kings[color] = square
		}
		file++
	}

	//to move
	if split[1] == "b" {
		to_move = 1
	}

	//castle rights
	for _, r := range split[2] {
		switch r {
		case 'K':
			castle_rights |= 0b0001
		case 'Q':
			castle_rights |= 0b0010
		case 'k':
			castle_rights |= 0b0100
		case 'q':
			castle_rights |= 0b1000
		}
	}
	//ep square
	if split[3] != "-" {
		ep_square = 8*(split[3][1]-'1') + split[3][0] - 'a'
	}
	var x uint64
	//full move
	x, _ = strconv.ParseUint(split[4], 10, 16)
	full_move = uint16(x)
	//half move
	x, _ = strconv.ParseUint(split[5], 10, 8)
	half_move = uint8(x)

	//derived bit boards
	for i := range 6 {
		allBB[0] |= pieceBB[i]
		allBB[1] |= pieceBB[i+6]
	}
	occupant = allBB[0] | allBB[1]

	return Position{pieceBB, allBB, occupant, to_move, castle_rights,
		ep_square, full_move, half_move, kings, moveStack, stateStack, ply}
}

// exportFen converts in-memory state back into a FEN string.
// TODO: when exporting mid make/unmake states dont make sense look into it
func (p *Position) exportFen() string {
	var sb strings.Builder
	var count uint8

	var board [64]uint8
	for p, bb := range p.pieceBB {
		for bb != 0 {
			board[popLSB(&bb)] = uint8(p) + 1
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
				sb.WriteByte(count + '0')
				count = 0
			}
			sb.WriteByte(Uint8ToChar[(c-1)%6] + (c-1)/6*32)
		}
		if count > 0 {
			sb.WriteByte(count + '0')
			count = 0
		}
		if rank != 0 {
			sb.WriteByte('/')
		}
	}

	if p.to_move == 0 {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	if p.castle_rights == 0 {
		sb.WriteString("-")
	} else {
		if p.castle_rights&0b0001 != 0 {
			sb.WriteByte('K')
		}
		if p.castle_rights&0b0010 != 0 {
			sb.WriteByte('Q')
		}
		if p.castle_rights&0b0100 != 0 {
			sb.WriteByte('k')
		}
		if p.castle_rights&0b1000 != 0 {
			sb.WriteByte('q')
		}
	}
	if p.ep_square == 64 {
		sb.WriteString(" - ")
	} else {
		sb.WriteByte(' ')
		sb.WriteByte('a' + p.ep_square%8)
		sb.WriteByte('1' + p.ep_square/8)
		sb.WriteByte(' ')
	}
	var buf [8]byte
	b := strconv.AppendUint(buf[:0], uint64(p.full_move), 10)
	b = append(b, ' ')
	b = strconv.AppendUint(b, uint64(p.half_move), 10)
	sb.Write(b)

	return sb.String()
}
