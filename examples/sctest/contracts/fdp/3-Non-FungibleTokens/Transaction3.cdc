// Transaction3.cdc

import NonFungibleToken from 0x02

transaction {
		prepare(acct: Account) {
				// get the owners collection reference
				let receiverRef = acct.published[&NonFungibleToken.NFTReceiver] ?? panic("No receiver")
				// use the factory reference to mint an NFT, which deposits
				// the NFT into the collection that is sent as a parameter
				acct.storage[NonFungibleToken.NFTMinter]?.mintNFT(recipient: receiverRef)
		}
		execute {}
}