// Transaction2.cdc

import FungibleToken from 0x01

transaction {
	prepare(acct: Account) {
		// create a new empty Vault object
		let vaultA <- FungibleToken.createEmptyVault()
			
		// store it in the account storage
		// and destroy whatever was there previously
		let oldVault <- acct.storage[FungibleToken.Vault] <- vaultA
		destroy oldVault

        log("Empty Vault stored")

		// publish a new Receiver reference to the Vault
		acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

        log("Reference created")
	}

    post {
        // make sure the reference was created correctly
        getAccount(0x02).published[&FungibleToken.Receiver] != nil:  "Vault Receiver Reference was not created correctly"
    }
}
 