access(all) contract FungibleToken {

    access(all) var totalSupply: UInt64

	// Interface that enforces the requirements for withdrawing
    // tokens from the implementing type
    //
    access(all) resource interface Provider {
        access(all) fun withdraw(amount: UInt64): @Vault {
            post {
                result.balance == amount:
                    "Withdrawal amount must be the same as the balance of the withdrawn Vault"
            }
        }
    }

	// Interface that enforces the requirements for depositing
    // tokens into the implementing type
    //
	access(all) resource interface Receiver {
		access(all) var balance: UInt64

        access(all) fun deposit(from: @Vault) {
            pre {
                from.balance > UInt64(0):
                    "Deposit balance must be positive"
            }
        }
    }

    access(all) resource Vault: Receiver {
        
		// keeps track of the total balance of the account's tokens
        access(all) var balance: UInt64

        init(balance: UInt64) {
            self.balance = balance
        }

		// withdraw subtracts amount from the vaults balance and 
        // returns a vault object with the subtracted balance
        access(all) fun withdraw(amount: UInt64): @Vault {
            self.balance = self.balance - amount
            return <-create Vault(balance: amount)
        }
        
		// deposit takes a vault object as a parameter and adds
        // its balance to the balance of the Account's vault, then
        // destroys the sent vault because its balance has been consumed
        access(all) fun deposit(from: @Vault) {
            self.balance = self.balance + from.balance
            destroy from
        }
    }

	// allows anyone to be able to create a new, empty Vault so that
	// they can send and receive tokens
    access(all) fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0)
    }

	// Resource object that an admin can control to mint new tokens
    access(all) resource VaultMinter {
		// function that mints new tokens and deposits into an account's vault
		// using their `Receiver` reference
        access(all) fun mintTokens(amount: UInt64, recipient: &Receiver) {
			FungibleToken.totalSupply = FungibleToken.totalSupply + UInt64(amount)
            recipient.deposit(from: <-create Vault(balance: amount))
        }
    }

    init() {
        self.totalSupply = 30

        // create the Vault with the initial balance and put it in storage
        let oldVault <- self.account.storage[Vault] <- create Vault(balance: 30)
        destroy oldVault

		let oldMinter <- self.account.storage[VaultMinter] <- create VaultMinter()
        destroy oldMinter
    }
}
 