// Transaction2.cdc

import FungibleToken from 0x01

// This transaction configures an account to store and receive tokens defined by
// the FungibleToken contract.
transaction {
	prepare(acct: AuthAccount) {
		// Create a new empty Vault object
		let vaultA <- FungibleToken.createEmptyVault()
			
		// Store the vault in the account storage
		// and destroy whatever was there previously
		let oldVault <- acct.storage[FungibleToken.Vault] <- vaultA
		destroy oldVault

        log("Empty Vault stored")

		// Publish a new Receiver reference to the Vault
		acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

        log("Reference created")
	}

    post {
        // Check that the reference was created correctly
        getAccount(0x02).published[&FungibleToken.Receiver] != nil:  "Vault Receiver Reference was not created correctly"
    }
}
