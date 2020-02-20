// Transaction4.cdc

import NonFungibleToken from 0x02

transaction {
		prepare(acct: Account) {

				// create a new empty collection
				let collection <- NonFungibleToken.createEmptyCollection()
			
				// put it in storage
				let oldCollection <- acct.storage[NonFungibleToken.Collection] <- collection
				destroy oldCollection

				// publish a public interface
				acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.NFTReceiver
		}
}
