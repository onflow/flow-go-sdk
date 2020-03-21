import FungibleToken from 0x01

// This transaction is a template for a transaction that
// could be used by anyone to send tokens to another account
// that owns a Vault
transaction {

    // Temporary Vault object that holds the balance that is being transferred
    var temporaryVault: @FungibleToken.Vault

    prepare(acct: AuthAccount) {
        // withdraw tokens from your vault
        self.temporaryVault <- acct.storage[FungibleToken.Vault]?.withdraw(amount: 10) ?? panic("No Vault!")
    }

    execute {
        // get the recipient's public account object
        let recipient = getAccount(0x01)

        // get the recipient's Receiver reference to their Vault
        let receiverRef = recipient.published[&FungibleToken.Receiver] ?? panic("No receiver!")

        // deposit your tokens to their Vault
        receiverRef.deposit(from: <-self.temporaryVault)

        log("Transfer succeeded!")
    }
}
