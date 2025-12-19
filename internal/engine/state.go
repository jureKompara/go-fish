package engine

// type representing a state prior to a move that was made
type State struct {
	capture      uint8
	castleRights uint8
	epSquare     uint8
	halfmove     uint8
	hash         uint64
	//add hash and other stuff that is nice to cache here
}
