package flow

import (
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

const (
	// AddressLength is the size of an account address.
	AddressLength = 20
)

func init() {
	gob.Register(Address{})
}

// An Address is a 20-byte identifier for a Flow account.
type Address [AddressLength]byte

var (
	// ZeroAddress represents the "zero address" (account that no one owns).
	ZeroAddress = Address{}
	// RootAddress is the address of the Flow root account.
	RootAddress = BytesToAddress(big.NewInt(1).Bytes())
)

// BytesToAddress returns Address with value b.
//
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// HexToAddress converts a hex string to an Address.
func HexToAddress(h string) Address {
	b, _ := hex.DecodeString(h)
	return BytesToAddress(b)
}

// SetBytes sets this address to the value of b.
//
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// Bytes returns the byte representation of this address.
func (a Address) Bytes() []byte { return a[:] }

// Hex returns the hex string representation of this address.
func (a Address) Hex() string {
	return hex.EncodeToString(a.Bytes())
}

// String returns the string representation of this address.
func (a Address) String() string {
	return a.Hex()
}

// Short returns the string representation of this address with leading zeros removed.
func (a Address) Short() string {
	hex := a.String()
	trimmed := strings.TrimLeft(hex, "0")
	if len(trimmed)%2 != 0 {
		trimmed = "0" + trimmed
	}
	return trimmed
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", a.Hex())), nil
}

func (a *Address) UnmarshalJSON(data []byte) error {
	*a = HexToAddress(strings.Trim(string(data), "\""))
	return nil
}
