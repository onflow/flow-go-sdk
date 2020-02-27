// Transaction4.cdc

import NonFungibleToken from 0x02

// This transaction is how a user would set up their account
// to use the NFT by creating a new empty collection,
// storing it in their account storage, and publishing a reference.
transaction {
    prepare(acct: Account) {

        // create a new empty collection
        let collection <- NonFungibleToken.createEmptyCollection()
    
        // put it in storage
        let oldCollection <- acct.storage[NonFungibleToken.Collection] <- collection
        destroy oldCollection

        // publish a public reference
        acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.NFTReceiver
    }
}
