// FungibleToken.cdc
//
// The FungibleToken contract is a sample implementation of a fungible token on Flow.
//
// Fungible tokens behave like everyday currencies -- they can be minted, transferred or
// traded for digital goods.
//
// Follow the fungible tokens tutorial to learn more: https://docs.onflow.org/docs/fungible-tokens

access(all) contract FungibleToken {

    // Total supply of all tokens in existence.
    access(all) var totalSupply: UInt64

    // Provider
    // 
    // Interface that enforces the requirements for withdrawing
    // tokens from the implementing type.
    //
    // We don't enforce requirements on self.balance here because
    // it leaves open the possibility of creating custom providers
    // that don't necessarily need their own balance.
    //
    access(all) resource interface Provider {

        // withdraw
        //
        // Function that subtracts tokens from the owner's Vault
        // and returns a Vault resource (@Vault) with the removed tokens.
        //
        // The function's access level is public, but this isn't a problem
        // because even the public functions are not fully public at first.
        // anyone in the network can call them, but only if the owner grants
        // them access by publishing a resource that exposes the withdraw
        // function.
        //
        access(all) fun withdraw(amount: UInt64): @Vault {
            post {
                // `result` refers to the return value of the function
                result.balance == amount:
                    "Withdrawal amount must be the same as the balance of the withdrawn Vault"
            }
        }
    }

    // Receiver 
    //
    // Interface that enforces the requirements for depositing
    // tokens into the implementing type.
    //
    // We don't include a condition that checks the balance because
    // we want to give users the ability to make custom Receivers that
    // can do custom things with the tokens, like split them up and
    // send them to different places.
    //
	access(all) resource interface Receiver {
		access(all) var balance: UInt64

        // deposit
        //
        // Function that can be called to deposit tokens 
        // into the implementing resource type
        //
        access(all) fun deposit(from: @Vault) {
            pre {
                from.balance > UInt64(0):
                    "Deposit balance must be positive"
            }
        }
    }

    // Vault
    //
    // Each user stores an instance of only the Vault in their storage
    // The functions in the Vault and governed by the pre and post conditions
    // in the interfaces when they are called. 
    // The checks happen at runtime whenever a function is called.
    //
    // Resources can only be created in the context of the contract that they
    // are defined in, so there is no way for a malicious user to create Vaults
    // out of thin air. A special Minter resource needs to be defined to mint
    // new tokens.
    // 
    access(all) resource Vault: Receiver {
        
		// keeps track of the total balance of the account's tokens
        access(all) var balance: UInt64

        // initialize the balance at resource creation time
        init(balance: UInt64) {
            self.balance = balance
        }

        // withdraw
        //
        // Function that takes an integer amount as an argument
        // and withdraws that amount from the Vault.
        //
        // It creates a new temporary Vault that is used to hold
        // the money that is being transferred. It returns the newly
        // created Vault to the context that called so it can be deposited
        // elsewhere.
        //
        access(all) fun withdraw(amount: UInt64): @Vault {
            self.balance = self.balance - amount
            return <-create Vault(balance: amount)
        }
        
        // deposit
        //
        // Function that takes a Vault object as an argument and adds
        // its balance to the balance of the owners Vault.
        //
        // It is allowed to destroy the sent Vault because the Vault
        // was a temporary holder of the tokens. The Vault's balance has
        // been consumed and therefore can be destroyed.
        access(all) fun deposit(from: @Vault) {
            self.balance = self.balance + from.balance
            destroy from
        }
    }

    // createEmptyVault
    //
    // Function that creates a new Vault with a balance of zero
    // and returns it to the calling context. A user must call this function
    // and store the returned Vault in their storage in order to allow their
    // account to be able to receive deposits of this token type.
    //
    access(all) fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0)
    }

	// VaultMinter
    //
    // Resource object that an admin can control to mint new tokens
    access(all) resource VaultMinter {

		// Function that mints new tokens and deposits into an account's vault
		// using their `Receiver` reference.
        access(all) fun mintTokens(amount: UInt64, recipient: &Receiver) {
			FungibleToken.totalSupply = FungibleToken.totalSupply + UInt64(amount)
            recipient.deposit(from: <-create Vault(balance: amount))
        }
    }

    // The init function for the contract. All fields in the contract must
    // be initialized at deployment. This is just an example of what
    // an implementation could do in the init function. The numbers are arbitrary.
    init() {
        self.totalSupply = 30

        // Create the Vault with the initial balance and put it into storage.
        let oldVault <- self.account.storage[Vault] <- create Vault(balance: 30)
        destroy oldVault

        // Create a VaultMinter resource object and put it into storage.
		let oldMinter <- self.account.storage[VaultMinter] <- create VaultMinter()
        destroy oldMinter
    }
}
