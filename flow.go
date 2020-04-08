package flow

import "encoding/hex"

// Identifier represents a 32-byte unique identifier for an entity.
type Identifier [32]byte

var ZeroID = Identifier{}

func (i Identifier) Bytes() []byte {
	return i[:]
}

func (i Identifier) Hex() string {
	return hex.EncodeToString(i[:])
}

func (i Identifier) String() string {
	return i.Hex()
}

func HashToID(hash []byte) Identifier {
	var id Identifier
	copy(id[:], hash)
	return id
}
