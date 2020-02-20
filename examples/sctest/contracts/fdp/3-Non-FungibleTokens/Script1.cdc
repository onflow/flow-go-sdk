// Script1.cdc 

import NonFungibleToken from 0x02

// Print Account NFTs
access(all) fun main() {
		let nftOwner = getAccount(0x02)

		let collectionRef = nftOwner.published[&NonFungibleToken.NFTReceiver] ?? panic("missing reference!")

		log("Account 2 NFTs")
		log(collectionRef.getIDs())
}