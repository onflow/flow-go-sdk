// SetupAccount1TransactionMinting.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// This transaction is what needs to be run from account 1
// after account 2 has been setup to finish minting tokens
// so that the accounts are ready for the Marketplace tutorial
transaction {

    // Public Vault Receiver References for both accounts
    let acct1Ref: &FungibleToken.Receiver
    let acct2Ref: &FungibleToken.Receiver

    // Private minter references for this account to mint tokens
    let minterRef: &FungibleToken.VaultMinter     
    
    prepare(acct: Account) {
        // get account objects
        let account2 = getAccount(0x02)
        
        // retreive public Vault Receiver references
        self.acct1Ref = acct.published[&FungibleToken.Receiver] ?? panic("no receiver 1 Ref")
        self.acct2Ref = account2.published[&FungibleToken.Receiver] ?? panic("no receiver 2 Ref")

        // get the stored moment references
        self.minterRef = &acct.storage[FungibleToken.VaultMinter] as &FungibleToken.VaultMinter
    }
    execute {

        // Mint tokens for both accounts
        self.minterRef.mintTokens(amount: 20, recipient: self.acct2Ref)
        self.minterRef.mintTokens(amount: 10, recipient: self.acct1Ref)

        log("Minted new fungible tokens for account 1 and 2")
    }
}
