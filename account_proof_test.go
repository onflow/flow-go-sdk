package flow

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountProofMsg(t *testing.T) {
	type testCase struct {
		addr           Address
		timestamp      int64
		appDomainTag   string
		expectedResult string
		expectedErr    error
	}

	for name, tc := range map[string]testCase{
		"with domain tag": {
			addr:           HexToAddress("ABC123DEF456"),
			timestamp:      int64(1632179933495),
			appDomainTag:   "FLOW-JS-SDK",
			expectedResult: "f1a0464c4f572d4a532d53444b000000000000000000000000000000000000000000880000abc123def45686017c05815137",
		},
		"without domain tag": {
			addr:           HexToAddress("ABC123DEF456"),
			timestamp:      int64(1632179933495),
			expectedResult: "d0880000abc123def45686017c05815137",
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Check the output of NewAccountProofMessage against a pre-generated message from the flow-js-sdk
			msg, err := NewAccountProofMessage(tc.addr, tc.timestamp, tc.appDomainTag)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, hex.EncodeToString(msg))
		})
	}
}
