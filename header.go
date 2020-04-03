package flow

type BlockHeader struct {
	ID       Identifier
	ParentID Identifier
	Height   uint64
}
