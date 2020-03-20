// Transaction3.cdc

import FungibleToken from 0x01

// This transaction performs a token transfer from one account to another using the
// FungibleToken contract.
transaction {

    // Local variable for storing the reference to the minter resource
    let mintingRef: &FungibleToken.VaultMinter

    // Local variable for storing the reference to the Vault of
    // the account that will receive the newly minted tokens
    var receiverRef: &FungibleToken.Receiver

    // The balance of the account before the minting happens
    // used for the post condition of the transaction.
    var beforeBalance: UInt64

	prepare(acct: AuthAccount) {
        // Create a reference to the stored, private minter resource
        self.mintingRef = &acct.storage[FungibleToken.VaultMinter] as &FungibleToken.VaultMinter

        // Get the public account object for account 0x02
        let recipient = getAccount(0x02)

        // Find their published Receiver reference and record their balance
        self.receiverRef = recipient.published[&FungibleToken.Receiver] ?? panic("No receiver reference!")
        self.beforeBalance = self.receiverRef.balance
	}

    execute {
        // Mint 30 tokens and deposit them into the recipient's Vault
        self.mintingRef.mintTokens(amount: 30, recipient: self.receiverRef)

        log("30 tokens minted and deposited to account 0x02")
    }

    post {
        // Make sure their account balance has been
        // increased by 30 
        self.receiverRef.balance == self.beforeBalance + UInt64(30): "Minting failed"
    }
}
