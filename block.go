package flow

// A Block is a set of state mutations applied to the Flow blockchain.
type Block struct {
	BlockHeader
	BlockPayload
}

// A BlockHeader is a summary of a full block.
type BlockHeader struct {
	ID       Identifier
	ParentID Identifier
	Height   uint64
}

// A BlockPayload is the full contents of a block.
//
// A payload contains the collection guarantees and seals for a block.
type BlockPayload struct {
	CollectionGuarantees []*CollectionGuarantee
	Seals                []*BlockSeal
}

// TODO: define block seal struct
type BlockSeal struct{}
