// Script2.cdc 

import NonFungibleToken from 0x02

// Print the NFTs owned by account 0x02.
access(all) fun main() {
    // Get the public account object for account 0x02
    let nftOwner = getAccount(0x02)

    // Find the public Receiver reference to their Collection
    let collectionRef = nftOwner.published[&NonFungibleToken.NFTReceiver] ?? panic("missing reference!")

    // Log the NFTs that they own as an array of IDs
    log("Account 2 NFTs")
    log(collectionRef.getIDs())
}
