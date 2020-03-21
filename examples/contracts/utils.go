package contracts

import (
	"testing"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/emulator"
	"github.com/dapperlabs/flow-go-sdk/utils/examples"
)

// setupUsersTokens sets up two accounts with 30 Fungible Tokens each
// and a NFT collection with 1 NFT each
func setupUsersTokens(
	t *testing.T,
	b *emulator.Blockchain,
	tokenAddr flow.Address,
	nftAddr flow.Address,
	signingKeys []flow.AccountPrivateKey,
	signingAddresses []flow.Address,
) {
	// add array of signers to transaction
	for i := 0; i < len(signingAddresses); i++ {
		tx := flow.Transaction{
			Script:         GenerateCreateTokenScript(tokenAddr, 30),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{signingAddresses[i]},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), signingKeys[i]}, []flow.Address{b.RootAccountAddress(), signingAddresses[i]}, false)

		// then deploy a NFT to the accounts
		tx = flow.Transaction{
			Script:         GenerateCreateNFTScript(nftAddr, i+1),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{signingAddresses[i]},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), signingKeys[i]}, []flow.Address{b.RootAccountAddress(), signingAddresses[i]}, false)
	}
}
