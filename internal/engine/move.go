package engine

import "fmt"

type Move uint32

const (
	DP      uint8 = 0b00001
	EP      uint8 = 0b00010
	KCASTLE uint8 = 0b00100
	QCASTLE uint8 = 0b01000
	ISCAP   uint8 = 0b10000
)

func NewMove(from, to, piece, promo, capture int, flags uint8) Move {
	return Move(uint32(from) |
		uint32(to)<<6 |
		uint32(piece)<<12 |
		uint32(promo)<<15 |
		uint32(capture)<<18 |
		uint32(flags)<<21)
}

// geters for packed uint32 Move
func (m Move) From() int    { return int(m & 0x3F) }
func (m Move) To() int      { return int((m >> 6) & 0x3F) }
func (m Move) Piece() int   { return int((m >> 12) & 0x7) }
func (m Move) Promo() int   { return int((m >> 15) & 0x7) }
func (m Move) Capture() int { return int((m >> 18) & 0x7) }
func (m Move) Flags() uint8 { return uint8((m >> 21) & 0x1F) }

// converts a move to UCI notation (e4e5 c7c8q)
func (m Move) Uci() string {
	from := m.From()
	to := m.To()
	promo := m.Promo()

	r := from >> 3
	f := from & 7
	tr := to >> 3
	tf := to & 7

	if promo != EMPTY {
		return fmt.Sprintf("%c%d%c%d%c", f+'a', r+1, tf+'a', tr+1, _pieceToChar[promo])
	}
	return fmt.Sprintf("%c%d%c%d", f+'a', r+1, tf+'a', tr+1)

}
