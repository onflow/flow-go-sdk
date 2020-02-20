// Transaction2.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02
import Marketplace from 0x03

transaction {

		prepare(acct: Account) {

				// get the read-only account storage of the seller
				let seller = getAccount(0x01)

				// get the references to the buyer's fungible token Vault and NFT Collection Receiver
				let collectionRef = acct.published[&NonFungibleToken.NFTReceiver] ?? panic("missing collection reference!")
				let vaultRef = &acct.storage[FungibleToken.Vault] as &FungibleToken.Vault
			
				// withdraw tokens from the buyers Vault
				let tokens <- vaultRef.withdraw(amount: 10)

				// get the reference to the seller's sale
				let saleRef = seller.published[&Marketplace.SalePublic] ?? panic("missing sale reference!")
		
				// purchase the NFT the the seller is selling, giving them the reference
				// to your NFT collection and giving them the tokens to buy it
				saleRef.purchase(tokenID: 1, recipient: collectionRef, buyTokens: <-tokens)
		}
}
