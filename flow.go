package flow

// Identifier represents a 32-byte unique identifier for an entity.
type Identifier [32]byte

var ZeroID = Identifier{}

func (i Identifier) Bytes() []byte {
	return i[:]
}

func HashToID(hash []byte) Identifier {
	var id Identifier
	copy(id[:], hash)
	return id
}
