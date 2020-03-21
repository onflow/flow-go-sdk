package templates_test

import (
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
)

func TestCreateAccount(t *testing.T) {
	publicKey := unittest.PublicKeyFixtures()[0]

	accountKey := flow.AccountPublicKey{
		PublicKey: publicKey,
		SignAlgo:  publicKey.Algorithm(),
		HashAlgo:  crypto.SHA3_256,
		Weight:    keys.PublicKeyWeightThreshold,
	}

	t.Run("without code", func(t *testing.T) {

		script, err := templates.CreateAccount([]flow.AccountPublicKey{accountKey}, []byte{})
		assert.NoError(t, err)

		expectedScript :=`
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
