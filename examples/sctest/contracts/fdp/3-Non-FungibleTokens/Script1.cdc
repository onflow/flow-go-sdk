// Script1.cdc 

import NonFungibleToken from 0x03

// Prints the NFTs that a specifc account owns
access(all) fun main() {
        // get the acccounts public account object
		let nftOwner = getAccount(0x03)

        // find their public Receiver reference to their Collection
		let collectionRef = nftOwner.published[&NonFungibleToken.NFTReceiver] ?? panic("missing reference!")

        // Log the NFTs that they own as an array of IDs
		log("Account 3 NFTs")
		log(collectionRef.getIDs())
}