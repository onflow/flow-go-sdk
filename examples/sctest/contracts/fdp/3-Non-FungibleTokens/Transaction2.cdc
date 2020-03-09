// Transaction2.cdc

import NonFungibleToken from 0x03

// transaction that published a public reference
// to the stored NFT Collection
transaction {
    prepare(acct: Account) {
        // publish a public interface that 
        // only exposes ownedNFTs, deposit, getIDs, and idExists
        acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.NFTReceiver

        log("Collection Reference created successfully")
    }
}