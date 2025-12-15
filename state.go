package main

// type representing a state prior to a move that was made
type State uint16

func (s State) Capture() uint8      { return uint8(s & 0x7) }
func (s State) CastleRights() uint8 { return uint8(s >> 3 & 0xF) }
func (s State) EPsquare() uint8     { return uint8(s >> 7 & 0x7F) }
