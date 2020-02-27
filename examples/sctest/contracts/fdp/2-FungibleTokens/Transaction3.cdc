// Transaction3.cdc

import FungibleToken from 0x01

transaction {

    // local variable for storing the reference to the minter resource
    let mintingRef: &FungibleToken.VaultMinter

    // local variable for storing the reference to the Vault of
    // the account that will receive the newly minted tokens
    var receiverRef: &FungibleToken.Receiver

    // The balance of the account before the minting happens
    // used for the post condition of the transaction.
    var beforeBalance: UInt64

	prepare(acct: Account) {
        // Create a reference to the stored, private minter resource
        self.mintingRef = &acct.storage[FungibleToken.VaultMinter] as &FungibleToken.VaultMinter

        // get account 0x02's public account object
        let recipient = getAccount(0x02)

        // find their published Reciever reference
        // and record their balance
        self.receiverRef = recipient.published[&FungibleToken.Receiver] ?? panic("No receiver reference!")
        self.beforeBalance = self.receiverRef.balance
	}

    execute {
        // Mint 30 tokens and deposit them into the recipient's Vault
        self.mintingRef.mintTokens(amount: 30, recipient: self.receiverRef)
    }

    post {
        // Make sure their account balance has been
        // increased by 30 
        self.receiverRef.balance == self.beforeBalance + UInt64(30): "30 Tokens minted and deposited to account 0x02"
    }
}
 