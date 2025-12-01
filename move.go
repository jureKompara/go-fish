package main

import "fmt"

type Move uint32

const (
	DP      uint8 = 0b00001
	EP      uint8 = 0b00010
	KCASTLE uint8 = 0b00100
	QCASTLE uint8 = 0b01000
	ISCAP   uint8 = 0b10000
	NOPROMO uint8 = 7
	NOCAP   uint8 = 7
)

func NewMove(from, to, piece, promo, capture, flags uint8) Move {
	return Move(uint32(from) |
		uint32(to)<<6 |
		uint32(piece)<<12 |
		uint32(promo)<<15 |
		uint32(capture)<<18 |
		uint32(flags)<<21)
}

// geters for pact uint32 Move
func (m Move) From() uint8    { return uint8(m & 0x3F) }
func (m Move) To() uint8      { return uint8((m >> 6) & 0x3F) }
func (m Move) Piece() uint8   { return uint8((m >> 12) & 0x7) }
func (m Move) Promo() uint8   { return uint8((m >> 15) & 0x7) }
func (m Move) Capture() uint8 { return uint8((m >> 18) & 0x7) }
func (m Move) Flags() uint8   { return uint8((m >> 21) & 0x1F) }

// converts a move to standard algebraic notation
func (m Move) San() string {
	from := m.From()
	to := m.To()

	r := from & 0b11111000
	f := from & 0b00000111
	tr := to & 0b11111000
	tf := to & 0b00000111

	return fmt.Sprintf("%c%d%c%d", f+'a', r+1, tf+'a', tr+1)
}
