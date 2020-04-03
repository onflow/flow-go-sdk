package flow

// Identifier represents a 32-byte unique identifier for an entity.
type Identifier [32]byte

func HashToID(hash []byte) Identifier {
	var id Identifier
	copy(id[:], hash)
	return id
}
