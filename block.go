package flow

type Block struct {
	BlockHeader
	BlockPayload
}

type BlockHeader struct {
	ID       Identifier
	ParentID Identifier
	Height   uint64
}

type BlockPayload struct {
	Guarantees []*CollectionGuarantee
	Seals      []*BlockSeal
}

// TODO: define block seal struct
type BlockSeal struct{}
