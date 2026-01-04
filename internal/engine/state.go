package engine

// type representing a state prior to a move that was made
type State struct {
	hash         uint64
	capture      uint8
	castleRights uint8
	epSquare     uint8
	halfmove     uint8
}
