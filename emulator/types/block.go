package types

import (
	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go/model/hash"

	"github.com/dapperlabs/flow-go-sdk"
)

// Block is a naive data structure used to represent blocks in the emulator.
type Block struct {
	Number            uint64
	PreviousBlockHash crypto.Hash
	TransactionHashes []crypto.Hash
}

// Hash returns the hash of this block.
func (b Block) Hash() crypto.Hash {
	return hash.DefaultHasher.ComputeHash(b.Encode())
}

func (b Block) Encode() []byte {
	temp := struct {
		Number            uint64
		PreviousBlockHash crypto.Hash
		TransactionHashes []crypto.Hash
	}{
		b.Number,
		b.PreviousBlockHash,
		b.TransactionHashes,
	}

	return flow.DefaultEncoder.MustEncode(&temp)
}

// GenesisBlock returns the genesis block for an emulated blockchain.
func GenesisBlock() Block {
	return Block{
		Number:            0,
		PreviousBlockHash: nil,
		TransactionHashes: make([]crypto.Hash, 0),
	}
}
