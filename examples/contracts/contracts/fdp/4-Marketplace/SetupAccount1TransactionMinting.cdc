// SetupAccount1TransactionMinting.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// This transaction mints tokens for both accounts using
// the minter stored on account 0x01.
transaction {

    // Public Vault Receiver References for both accounts
    let acct1Ref: &FungibleToken.Receiver
    let acct2Ref: &FungibleToken.Receiver

    // Private minter references for this account to mint tokens
    let minterRef: &FungibleToken.VaultMinter     
    
    prepare(acct: AuthAccount) {
        // Get the public object for account 0x02
        let account2 = getAccount(0x02)
        
        // Retrieve public Vault Receiver references for both accounts
        self.acct1Ref = acct.published[&FungibleToken.Receiver] ?? panic("no receiver 1 Ref")
        self.acct2Ref = account2.published[&FungibleToken.Receiver] ?? panic("no receiver 2 Ref")

        // Get the stored Minter reference for account 0x01
        self.minterRef = &acct.storage[FungibleToken.VaultMinter] as &FungibleToken.VaultMinter
    }

    execute {
        // Mint tokens for both accounts
        self.minterRef.mintTokens(amount: 20, recipient: self.acct2Ref)
        self.minterRef.mintTokens(amount: 10, recipient: self.acct1Ref)

        log("Minted new fungible tokens for account 1 and 2")
    }
}
