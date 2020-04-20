package flow

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/onflow/flow-go-sdk/crypto"
)

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

// DefaultHasher is the default hasher used by Flow.
var DefaultHasher crypto.Hasher

func init() {
	DefaultHasher = crypto.NewSHA3_256()
}

func rlpEncode(v interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(v)
}

func rlpDecode(b []byte, v interface{}) error {
	return rlp.DecodeBytes(b, v)
}

func mustRLPEncode(v interface{}) []byte {
	b, err := rlpEncode(v)
	if err != nil {
		panic(err)
	}
	return b
}

func mustRLPDecode(b []byte, v interface{}) {
	err := rlpDecode(b, v)
	if err != nil {
		panic(err)
	}
}
