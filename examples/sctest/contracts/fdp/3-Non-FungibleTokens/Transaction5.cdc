// Transaction5.cdc

import NonFungibleToken from 0x02

transaction {
	
		prepare(acct: Account) {
				let recipient = getAccount(0x01)

				// get the Collection references for the receiver
				let depositRef = recipient.published[&NonFungibleToken.NFTReceiver] ?? panic("missing deposit reference")
		
				// call the withdraw function on the sender's Collection
				// to move the NFT out of the collection
				let token <- acct.storage[NonFungibleToken.Collection]?.withdraw(withdrawID: 1) ?? panic("missing collection")
		
				// deposit the NFT in the receivers collection
				depositRef.deposit(token: <-token)
		}
}
