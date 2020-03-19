// SetupAccount1Transaction.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// This transaction setting up account 1's account
// so that it is ready to use the marketplace tutorial
// by publishing a Vault reference and creating an empty NFT Collection
//
transaction {
        prepare(acct: Account) {
            // Create a public Receiver reference to the Vault
            let receiverRef = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver
            acct.published[&FungibleToken.Receiver] = receiverRef

            log("Created Vault references")

            // Create a new empty NFT Collection
            let collection <- NonFungibleToken.createEmptyCollection()
            
            // Put the NFT Collection in storage
            let oldCollection <- acct.storage[NonFungibleToken.Collection] <- collection
            destroy oldCollection

            // Publish a public interface to the Collection
            acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.Collection

            log("Created a new empty collection and published a reference")
        
        }
        execute {}
}
 