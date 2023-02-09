package crypto

import (
	"encoding/asn1"
	"fmt"
	"math/big"

)

// parseSignature parses an asn1 stucture (R,S) into a slice of bytes as required by the `Siger.Sign` method.
func ParseSignature(kmsSignature []byte, curve SignatureAlgorithm) ([]byte, error) {
	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(kmsSignature, &parsedSig); err != nil {
		return nil, fmt.Errorf("asn1.Unmarshal: %w", err)
	}

	curveOrderLen := curveOrder(curve)
	signature := make([]byte, 2*curveOrderLen)

	// left pad R and S with zeroes
	rBytes := parsedSig.R.Bytes()
	sBytes := parsedSig.S.Bytes()
	copy(signature[curveOrderLen-len(rBytes):], rBytes)
	copy(signature[len(signature)-len(sBytes):], sBytes)

	return signature, nil
}

// returns the curve order size in bytes (used to padd R and S of the ECDSA signature)
// Only P-256 and secp256k1 are supported. The calling function should make sure
// the function is only called with one of the 2 curves.
func curveOrder(curve SignatureAlgorithm) int {
	switch curve {
	case ECDSA_P256:
		return 32
	case ECDSA_secp256k1:
		return 32
	default:
		panic("the curve is not supported")
	}
}