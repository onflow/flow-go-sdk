// Transaction5.cdc

import NonFungibleToken from 0x02

// This transaction transfers an NFT from one user's collection
// to another user's collection.
transaction {

    // The field that will hold the NFT as it is being
    // transferred to the other account
    let transferToken: @NonFungibleToken.NFT
	
    prepare(acct: AuthAccount) {

        // Call the withdraw function on the sender's Collection
        // to move the NFT out of the collection
        self.transferToken <- acct.storage[NonFungibleToken.Collection]?.withdraw(withdrawID: 1) ?? panic("missing collection")
    }

    execute {
        // Get the recipient's public account object
        let recipient = getAccount(0x01)

        // Get the Collection reference for the receiver
        let receiverRef = recipient.published[&NonFungibleToken.NFTReceiver] ?? panic("missing deposit reference")

        // Deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-self.transferToken)

        log("NFT ID 1 transferred from account 2 to account 1")
    }
}
