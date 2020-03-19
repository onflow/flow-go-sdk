// Script1.cdc 

import NonFungibleToken from 0x03

// Print the NFTs owned by account 0x03.
access(all) fun main() {
    // Get the public account object for account 0x03
    let nftOwner = getAccount(0x03)

    // Find the public Receiver reference to their Collection
    let collectionRef = nftOwner.published[&NonFungibleToken.NFTReceiver] ?? panic("missing reference!")

    // Log the NFTs that they own as an array of IDs
    log("Account 3 NFTs")
    log(collectionRef.getIDs())
}
