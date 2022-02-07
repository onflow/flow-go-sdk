package flow

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountProofMsg(t *testing.T) {
	type testCase struct {
		addr           string
		timestamp      int64
		appDomainTag   string
		expectedResult string
		expectedErr    error
	}

	for name, tc := range map[string]testCase{
		"with domain tag": {
			addr:           "1234123412341234",
			timestamp:      int64(1644259653161),
			appDomainTag:   "AWESOME-APP-V0.0-user",
			expectedResult: "f852b8403431353734353533346634643435326434313530353032643536333032653330326437353733363537323030303030303030303030303030303030303030303088123412341234123486017ed5833629",
		},
		"without domain tag": {
			addr:           "1234123412341234",
			timestamp:      int64(1644259653161),
			expectedResult: "d088123412341234123486017ed5833629",
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Check the output of NewAccountProofMsg against a pre-generated message from the flow-js-sdk
			msg, err := NewAccountProofMsg(tc.addr, tc.timestamp, tc.appDomainTag)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, hex.EncodeToString(msg))
		})
	}
}
