package templates_test

import (
	"encoding/hex"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
)

const mockPublicKeyHex = "9347be39ddcfc2edc36381f5e7a86151c4055064f0e54cc472c36553655693e52a6cf801fabc17e399f36785b9e81c0d9416760a0739a9e7135d0c9d81a8ef5c"

var mockPublicKey crypto.PublicKey

func init() {
	bytesKey, _ := hex.DecodeString(mockPublicKeyHex)
	mockPublicKey, _ = crypto.DecodePublicKey(crypto.ECDSA_P256, bytesKey)
}

func TestCreateAccount(t *testing.T) {
	accountKey := flow.AccountKey{
		PublicKey: mockPublicKey,
		SignAlgo:  mockPublicKey.Algorithm(),
		HashAlgo:  crypto.SHA3_256,
		Weight:    flow.AccountKeyWeightThreshold,
	}

	t.Run("without code", func(t *testing.T) {
		script, err := templates.CreateAccount([]flow.AccountKey{accountKey}, []byte{})
		assert.NoError(t, err)

		expectedScript := `
          transaction {
            execute {
              AuthAccount(publicKeys: [[248,71,184,64,147,71,190,57,221,207,194,237,195,99,129,245,231,168,97,81,196,5,80,100,240,229,76,196,114,195,101,83,101,86,147,229,42,108,248,1,250,188,23,227,153,243,103,133,185,232,28,13,148,22,118,10,7,57,169,231,19,93,12,157,129,168,239,92,2,3,130,3,232]], code: "".decodeHex())
            }
          }
        `

		assert.Equal(t,
			dedent.Dedent(expectedScript),
			dedent.Dedent(string(script)),
		)
	})

	t.Run("with code", func(t *testing.T) {
		script, err := templates.CreateAccount([]flow.AccountKey{accountKey}, []byte("pub fun main() {}"))
		assert.Nil(t, err)

		expectedScript := `
          transaction {
            execute {
              AuthAccount(publicKeys: [[248,71,184,64,147,71,190,57,221,207,194,237,195,99,129,245,231,168,97,81,196,5,80,100,240,229,76,196,114,195,101,83,101,86,147,229,42,108,248,1,250,188,23,227,153,243,103,133,185,232,28,13,148,22,118,10,7,57,169,231,19,93,12,157,129,168,239,92,2,3,130,3,232]], code: "7075622066756e206d61696e2829207b7d".decodeHex())
            }
          }
        `

		assert.Equal(t,
			dedent.Dedent(expectedScript),
			dedent.Dedent(string(script)),
		)
	})
}

func TestUpdateAccountCode(t *testing.T) {
	script := templates.UpdateAccountCode([]byte("pub fun main() {}"))

	expectedScript := `
      transaction {
        prepare(signer: AuthAccount) {
          signer.setCode("7075622066756e206d61696e2829207b7d".decodeHex())
        }
      }
    `

	assert.Equal(t,
		dedent.Dedent(expectedScript),
		dedent.Dedent(string(script)),
	)
}
