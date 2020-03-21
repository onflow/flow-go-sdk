// Transaction2.cdc

import NonFungibleToken from 0x03

// This transaction publishes a public reference to the stored NFT collection.
transaction {
    prepare(acct: AuthAccount) {

        // Publish a public interface that only exposes "ownedNFTs", "deposit", "getIDs", and "idExists".
        acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.NFTReceiver

        log("Collection Reference created successfully")
    }
}
