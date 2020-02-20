// Transaction1.cdc

import NonFungibleToken from 0x02

transaction {
		prepare(acct: Account) {
				if acct.storage[NonFungibleToken.NFT] != nil {
						log("The token exists!")
				} else {
						log("No token found!")
				}
	  }
}