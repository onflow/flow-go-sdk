package flow

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

type canonicalAccountProofWithoutTag struct {
	Addr      []byte
	Timestamp uint64
}

type canonicalAccountProofWithTag struct {
	DomainTag []byte
	Addr      []byte
	Timestamp uint64
}

// NewAccountProofMessage creates a new account proof message for singing. The appDomainTag is optional and can be left
// empty. Note that the resulting byte slice does not contain the user domain tag.
func NewAccountProofMessage(address Address, timestamp int64, appDomainTag string) ([]byte, error) {
	var (
		encodedMsg []byte
		err        error
	)

	if appDomainTag != "" {
		paddedTag, err := NewDomainTag(appDomainTag)
		if err != nil {
			return nil, fmt.Errorf("error encoding domain tag: %w", err)
		}

		encodedMsg, err = rlp.EncodeToBytes(&canonicalAccountProofWithTag{
			DomainTag: paddedTag[:],
			Addr:      address.Bytes(),
			Timestamp: uint64(timestamp),
		})
	} else {
		encodedMsg, err = rlp.EncodeToBytes(&canonicalAccountProofWithoutTag{
			Addr:      address.Bytes(),
			Timestamp: uint64(timestamp),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("error encoding account proof message: %w", err)
	}

	return encodedMsg, nil
}

// NewDomainTag returns a new padded domain tag from the given string. This function returns an error if the domain
// tag is too long.
func NewDomainTag(tag string) (paddedTag [domainTagLength]byte, err error) {
	if len(tag) > domainTagLength {
		return paddedTag, fmt.Errorf("domain tag %s cannot be longer than %d characters", tag, domainTagLength)
	}

	return paddedDomainTag(tag), nil
}
