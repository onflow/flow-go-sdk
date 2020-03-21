// Transaction1.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02
import Marketplace from 0x03

// This transaction creates a new Sale Collection object,
// lists an NFT for sale, puts it in account storage,
// and creates a public reference to the sale so that others can buy the token.
transaction {

    prepare(acct: AuthAccount) {

        // Get a reference to the signer's vault Receiver
        let ownerVault = acct.published[&FungibleToken.Receiver] ?? panic("No receiver reference!")

        // Create a new Sale object, 
        // initializing it with the reference to the owner's vault
        let sale: @Marketplace.SaleCollection <- Marketplace.createSaleCollection(ownerVault: ownerVault)
    
        // Withdraw the NFT from their collection that they want to sell
        // and move it into the transaction's context
        let token <- acct.storage[NonFungibleToken.Collection]?.withdraw(withdrawID: 1) ?? panic("missing token!")

        // List it for sale by moving it into the sale object
        sale.listForSale(token: <-token, price: 10)

        // Store the sale object in the account storage 
        let oldSale <- acct.storage[Marketplace.SaleCollection] <- sale
        destroy oldSale

        // Create a public reference to the sale so that others can call its methods
        acct.published[&Marketplace.SalePublic] = &acct.storage[Marketplace.SaleCollection] as &Marketplace.SalePublic

        log("Sale Created for account 1. Selling NFT 1 for 10 tokens")
    }
}
