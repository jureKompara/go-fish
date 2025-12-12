package engine

import (
	"fmt"
)

type Move uint16

const (
	//flags < 2 are castles
	KCASTLE uint8 = iota
	QCASTLE
	QUIET
	DOUBLE

	CAPTURE
	EP
	FREE1
	FREE2
	//flags >= 8 are promotions
	PROMOKNIGHT
	PROMOBISHOP
	PROMOROOOK
	PROMOQUEEN
	//flags with 3rd bit set are captures
	PROMOKNIGHTX
	PROMOBISHOPX
	PROMOROOOKX
	PROMOQUEENX
)

func NewMove(from, to int, flags uint8) Move {
	return Move(uint16(from) |
		uint16(to)<<6 |
		uint16(flags)<<12)
}

// ///////////////////////////////12-15//6-11//0-5
// geters for packed uint16 Move[[flags][to][from]]
func (m Move) From() int    { return int(m & 0x3F) }
func (m Move) To() int      { return int((m >> 6) & 0x3F) }
func (m Move) Flags() uint8 { return uint8((m >> 12) & 0xF) }

func IsCapture(flag uint8) bool { return flag&4 != 0 }
func IsEP(flag uint8) bool      { return flag == EP }
func IsDP(flag uint8) bool      { return flag == DOUBLE }
func IsPromo(flag uint8) bool   { return flag >= 8 }
func IsCastle(flag uint8) bool  { return flag <= 1 }

// this only makes sense if the move is a promo
func Promo(flag uint8) uint8 {
	return flag&3 + KNIGHT
}

// converts a move to UCI notation (e4e5 c7c8q)
func (m Move) Uci() string {
	from := m.From()
	to := m.To()
	flag := m.Flags()

	r := from >> 3
	f := from & 7
	tr := to >> 3
	tf := to & 7

	//TODO: MAKE THIS NOT SUCK ASS
	if flag == PROMOKNIGHT || flag == PROMOKNIGHTX {
		return fmt.Sprintf("%c%d%c%dk", f+'a', r+1, tf+'a', tr+1)
	}
	if flag == PROMOBISHOP || flag == PROMOBISHOPX {
		return fmt.Sprintf("%c%d%c%db", f+'a', r+1, tf+'a', tr+1)
	}
	if flag == PROMOROOOK || flag == PROMOROOOKX {
		return fmt.Sprintf("%c%d%c%dr", f+'a', r+1, tf+'a', tr+1)
	}
	if flag == PROMOQUEEN || flag == PROMOQUEENX {
		return fmt.Sprintf("%c%d%c%dq", f+'a', r+1, tf+'a', tr+1)
	}
	return fmt.Sprintf("%c%d%c%d", f+'a', r+1, tf+'a', tr+1)

}
