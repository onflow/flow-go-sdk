// Transaction4.cdc

import NonFungibleToken from 0x02

// This transaction configures a user's account
// to use the NFT contract by creating a new empty collection,
// storing it in their account storage, and publishing a reference.
transaction {
    prepare(acct: AuthAccount) {

        // Create a new empty collection
        let collection <- NonFungibleToken.createEmptyCollection()
    
        // Put it in storage
        let oldCollection <- acct.storage[NonFungibleToken.Collection] <- collection
        destroy oldCollection

        log("Collection created for account 1")

        // Publish a public reference
        acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.NFTReceiver

        log("Reference published")
    }
}
