package contracts

import (
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
)

// GenerateCreateTokenScript creates a script that instantiates
// a new Vault instance and stores it in memory.
// balance is an argument to the Vault constructor.
// The Vault must have been deployed already.
func GenerateCreateTokenScript(tokenAddr flow.Address, initialBalance int) []byte {
	template := `
      import FungibleToken, FlowToken from 0x%s

      transaction {

          prepare(acct: AuthAccount) {
              let oldVault <- acct.storage[FlowToken.Vault] <- FlowToken.createVault(initialBalance: %d)
              destroy oldVault

              acct.published[&AnyResource{FungibleToken.Receiver}] =
                  &acct.storage[FlowToken.Vault] as &AnyResource{FungibleToken.Receiver}

              acct.published[&AnyResource{FungibleToken.Provider}] =
                  &acct.storage[FlowToken.Vault] as &AnyResource{FungibleToken.Provider}

              acct.published[&AnyResource{FungibleToken.Balance}] =
                  &acct.storage[FlowToken.Vault] as &AnyResource{FungibleToken.Balance}
          }
      }
    `
	return []byte(fmt.Sprintf(template, tokenAddr, initialBalance))
}

// GenerateCreateThreeTokensArrayScript creates a script
// that creates three new vault instances, stores them
// in an array of vaults, and then stores the array
// to the storage of the signer's account
func GenerateCreateThreeTokensArrayScript(tokenAddr flow.Address, initialBalance int, bal2 int, bal3 int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s

		transaction {

		  prepare(acct: AuthAccount) {
			let vaultA <- FlowToken.createVault(initialBalance: %d)
    		let vaultB <- FlowToken.createVault(initialBalance: %d)
			let vaultC <- FlowToken.createVault(initialBalance: %d)
			
			var vaultArray <- [<-vaultA, <-vaultB]

			vaultArray.append(<-vaultC)

			let storedVaults <- acct.storage[[FlowToken.Vault]] <- vaultArray
			destroy storedVaults

            acct.published[&[FlowToken.Vault]] = &acct.storage[[FlowToken.Vault]] as &[FlowToken.Vault]
		  }
		}
	`
	return []byte(fmt.Sprintf(template, tokenAddr, initialBalance, bal2, bal3))
}

// GenerateWithdrawScript creates a script that withdraws
// tokens from a vault and destroys the tokens
func GenerateWithdrawScript(tokenCodeAddr flow.Address, vaultNumber int, withdrawAmount int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s

		transaction {
		  prepare(acct: AuthAccount) {
			var vaultArray <- acct.storage[[FlowToken.Vault]] ?? panic("missing vault array!")
			
			let withdrawVault <- vaultArray[%d].withdraw(amount: %d)

			var storedVaults: @[FlowToken.Vault]? <- vaultArray
			acct.storage[[FlowToken.Vault]] <-> storedVaults

			destroy withdrawVault
			destroy storedVaults
		  }
		}
	`

	return []byte(fmt.Sprintf(template, tokenCodeAddr, vaultNumber, withdrawAmount))
}

// GenerateWithdrawDepositScript creates a script
// that withdraws tokens from a vault and deposits
// them to another vault
func GenerateWithdrawDepositScript(tokenCodeAddr flow.Address, withdrawVaultNumber int, depositVaultNumber int, withdrawAmount int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s

		transaction {
		  prepare(acct: AuthAccount) {
			var vaultArray <- acct.storage[[FlowToken.Vault]] ?? panic("missing vault array!")
			
			let withdrawVault <- vaultArray[%d].withdraw(amount: %d)

			vaultArray[%d].deposit(from: <-withdrawVault)

			var storedVaults: @[FlowToken.Vault]? <- vaultArray
			acct.storage[[FlowToken.Vault]] <-> storedVaults

			destroy storedVaults
		  }
		}
	`

	return []byte(fmt.Sprintf(template, tokenCodeAddr, withdrawVaultNumber, withdrawAmount, depositVaultNumber))
}

// GenerateDepositVaultScript creates a script that withdraws an tokens from an account
// and deposits it to another account's vault
func GenerateDepositVaultScript(tokenCodeAddr flow.Address, receiverAddr flow.Address, amount int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s

		transaction {
		  prepare(acct: AuthAccount) {
			let recipient = getAccount(0x%s)

			let providerRef = acct.published[&AnyResource{FungibleToken.Provider}] ?? panic("missing Provider reference")
			let receiverRef = recipient.published[&AnyResource{FungibleToken.Receiver}] ?? panic("missing Receiver reference")

			let tokens <- providerRef.withdraw(amount: %d)

			receiverRef.deposit(from: <-tokens)
		  }
		}
	`

	return []byte(fmt.Sprintf(template, tokenCodeAddr, receiverAddr, amount))
}

// GenerateInspectVaultScript creates a script that retrieves a
// Vault from the array in storage and makes assertions about
// its balance. If these assertions fail, the script panics.
func GenerateInspectVaultScript(tokenCodeAddr, userAddr flow.Address, expectedBalance int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let vaultRef = acct.published[&AnyResource{FungibleToken.Balance}] ?? panic("missing Receiver reference")
			assert(
                vaultRef.balance == UInt64(%d),
                message: "incorrect Balance!"
            )
		}
    `

	return []byte(fmt.Sprintf(template, tokenCodeAddr, userAddr, expectedBalance))
}

// GenerateInspectVaultArrayScript creates a script that retrieves a
// Vault from the array in storage and makes assertions about
// its balance. If these assertions fail, the script panics.
func GenerateInspectVaultArrayScript(tokenCodeAddr, userAddr flow.Address, vaultNumber int, expectedBalance int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let vaultArray = acct.published[&[FlowToken.Vault]] ?? panic("missing vault")
			assert(
                vaultArray[%d].balance == UInt64(%d),
                message: "incorrect Balance!"
            )
        }
	`

	return []byte(fmt.Sprintf(template, tokenCodeAddr, userAddr, vaultNumber, expectedBalance))
}
