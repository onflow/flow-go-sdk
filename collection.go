package flow

// A Collection is a list of transactions bundled together for inclusion in a block.
type Collection struct {
	TransactionIDs []Identifier
}

// ID returns the canonical SHA3-256 hash of this collection.
func (c Collection) ID() Identifier {
	return HashToID(DefaultHasher.ComputeHash(c.Encode()))
}

// Encode returns the canonical encoding of this collection.
func (c Collection) Encode() []byte {
	transactionIDs := make([][]byte, len(c.TransactionIDs))
	for i, id := range c.TransactionIDs {
		transactionIDs[i] = id.Bytes()
	}

	temp := struct {
		TransactionIDS [][]byte
	}{
		TransactionIDS: transactionIDs,
	}
	return mustRLPEncode(&temp)
}

// A CollectionGuarantee is an attestation signed by the nodes that have guaranteed a collection.
type CollectionGuarantee struct {
	CollectionID Identifier
}
