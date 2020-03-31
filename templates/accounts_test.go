package templates_test

import (
	"encoding/hex"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

const mockPublicKeyHex = "3059301306072a8648ce3d020106082a8648ce3d0301070342000472b074a452d0a764a1da34318f44cb16740df1cfab1e6b50e5e4145dc06e5d151c9c25244f123e53c9b6fe237504a37e7779900aad53ca26e3b57c5c3d7030c4"

var mockPublicKey crypto.PublicKey

func init() {
	bytesKey, _ := hex.DecodeString(mockPublicKeyHex)
	mockPublicKey, _ = crypto.DecodePublicKey(crypto.ECDSA_P256, bytesKey)
}

func TestCreateAccount(t *testing.T) {
	accountKey := flow.AccountPublicKey{
		PublicKey: mockPublicKey,
		SignAlgo:  mockPublicKey.Algorithm(),
		HashAlgo:  crypto.SHA3_256,
		Weight:    keys.PublicKeyWeightThreshold,
	}

	t.Run("without code", func(t *testing.T) {
		script, err := templates.CreateAccount([]flow.AccountPublicKey{accountKey}, []byte{})
		assert.NoError(t, err)

		expectedScript := `
          transaction {
            execute {
              AuthAccount(publicKeys: [[248,98,184,91,48,89,48,19,6,7,42,134,72,206,61,2,1,6,8,42,134,72,206,61,3,1,7,3,66,0,4,114,176,116,164,82,208,167,100,161,218,52,49,143,68,203,22,116,13,241,207,171,30,107,80,229,228,20,93,192,110,93,21,28,156,37,36,79,18,62,83,201,182,254,35,117,4,163,126,119,121,144,10,173,83,202,38,227,181,124,92,61,112,48,196,2,3,130,3,232]], code: "".decodeHex())
            }
          }
        `

		assert.Equal(t,
			dedent.Dedent(expectedScript),
			dedent.Dedent(string(script)),
		)
	})

	t.Run("with code", func(t *testing.T) {
		script, err := templates.CreateAccount([]flow.AccountPublicKey{accountKey}, []byte("pub fun main() {}"))
		assert.Nil(t, err)

		expectedScript := `
          transaction {
            execute {
              AuthAccount(publicKeys: [[248,98,184,91,48,89,48,19,6,7,42,134,72,206,61,2,1,6,8,42,134,72,206,61,3,1,7,3,66,0,4,114,176,116,164,82,208,167,100,161,218,52,49,143,68,203,22,116,13,241,207,171,30,107,80,229,228,20,93,192,110,93,21,28,156,37,36,79,18,62,83,201,182,254,35,117,4,163,126,119,121,144,10,173,83,202,38,227,181,124,92,61,112,48,196,2,3,130,3,232]], code: "7075622066756e206d61696e2829207b7d".decodeHex())
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
