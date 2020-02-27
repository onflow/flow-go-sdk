// SetupAccount1Transaction.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// This transaction is for setting up account 1's account
// so that it is ready to use the marketplace tutorial
transaction {
        prepare(acct: Account) {
            // create reference to the Vault
            let receiverRef = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver
            acct.published[&FungibleToken.Receiver] = receiverRef

            // create a new empty collection
            let collection <- NonFungibleToken.createEmptyCollection()
            
            // put it in storage
            let oldCollection <- acct.storage[NonFungibleToken.Collection] <- collection
            destroy oldCollection

            // publish a public interface that only exposes ownedNFTs and deposit
            acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.Collection
        
        }
        execute {}
}
 