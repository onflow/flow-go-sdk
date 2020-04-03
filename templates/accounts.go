package templates

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"
)

// CreateAccount generates a script that creates a new account.
func CreateAccount(accountKeys []flow.AccountKey, code []byte) ([]byte, error) {
	publicKeys := make([][]byte, len(accountKeys))

	for i, accountKey := range accountKeys {
		accountKeyBytes, err := keys.EncodePublicKey(accountKey)
		if err != nil {
			return nil, err
		}

		publicKeys[i] = accountKeyBytes
	}

	publicKeysStr := languageEncodeBytesArray(publicKeys)

	script := fmt.Sprintf(
		`
            transaction {
              execute {
                AuthAccount(publicKeys: %s, code: "%s".decodeHex())
              }
            }
        `,
		publicKeysStr,
		hex.EncodeToString(code),
	)

	return []byte(script), nil
}

// UpdateAccountCode generates a script that updates the code associated with an account.
func UpdateAccountCode(code []byte) []byte {
	return []byte(fmt.Sprintf(
		`
            transaction {
              prepare(signer: AuthAccount) {
                signer.setCode("%s".decodeHex())
              }
            }
        `,
		hex.EncodeToString(code),
	))
}

// AddAccountKey generates a script that adds a key to an account.
func AddAccountKey(accountKey flow.AccountKey) ([]byte, error) {
	accountKeyBytes, err := keys.EncodePublicKey(accountKey)
	if err != nil {
		return nil, err
	}

	publicKeyStr := languageEncodeBytes(accountKeyBytes)

	script := fmt.Sprintf(`
        transaction {
          prepare(signer: AuthAccount) {
            signer.addPublicKey(%s)
          }
        }
   	`, publicKeyStr)

	return []byte(script), nil
}

// RemoveAccountKey generates a script that removes a key from an account.
func RemoveAccountKey(index int) []byte {
	script := fmt.Sprintf(`
        transaction {
          prepare(signer: AuthAccount) {
            signer.removePublicKey(%d)
          }
        }
    `, index)

	return []byte(script)
}

// languageEncodeBytes converts a byte slice to a comma-separated list of uint8 integers.
func languageEncodeBytes(b []byte) string {
	if len(b) == 0 {
		return "[]"
	}

	return strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
}

// languageEncodeBytesArray converts a slice of byte slices to a comma-separated list of uint8 integers.
//
// Example: [][]byte{[]byte{1}, []byte{2,3}} -> "[[1],[2,3]]"
func languageEncodeBytesArray(b [][]byte) string {
	if len(b) == 0 {
		return "[]"
	}

	return strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
}

// languageEncodeIntArray converts a slice of integers to a comma-separated list.
func languageEncodeIntArray(i []int) string {
	if len(i) == 0 {
		return "[]"
	}

	return strings.Join(strings.Fields(fmt.Sprintf("%d", i)), ",")
}
