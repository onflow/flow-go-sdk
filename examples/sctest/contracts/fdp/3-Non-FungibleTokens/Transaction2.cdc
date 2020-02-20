// Transaction2.cdc

import NonFungibleToken from 0x02

transaction {
		prepare(acct: Account) {
				// publish a public interface that 
				// only exposes ownedNFTs, deposit, getIDs, and idExists
				acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.NFTReceiver
		}
}