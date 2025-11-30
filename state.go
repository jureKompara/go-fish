package main

// type representing a state prior to a move that was made
type State struct {
	castleRights uint8
	epSquare     uint8
	halfmove     uint8
	//add hash and other stuff that is nice to cache here
}

func NewState(cr, ep, hm, cap uint8) State {
	return State{cr, ep, hm}
}
