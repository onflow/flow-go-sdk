// Transaction2.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02
import Marketplace from 0x03

// This transaction uses the signers Vault tokens to purchase an NFT
// from the Sale collection of account 0x01.
transaction {

    // reference to the buyer's NFT collection where they
    // will store the bought NFT
    let collectionRef: &NonFungibleToken.NFTReceiver

    // Vault that will hold the tokens that will be used to
    // but the NFT
    let temporaryVault: @FungibleToken.Vault

    prepare(acct: AuthAccount) {

        // get the references to the buyer's fungible token Vault and NFT Collection Receiver
        self.collectionRef = acct.published[&NonFungibleToken.NFTReceiver] ?? panic("missing collection reference!")
        let vaultRef = &acct.storage[FungibleToken.Vault] as &FungibleToken.Vault
    
        // withdraw tokens from the buyers Vault
        self.temporaryVault <- vaultRef.withdraw(amount: 10)
    }

    execute {
        // get the read-only account storage of the seller
        let seller = getAccount(0x01)

        // get the reference to the seller's sale
        let saleRef = seller.published[&Marketplace.SalePublic] ?? panic("missing sale reference!")

        // purchase the NFT the the seller is selling, giving them the reference
        // to your NFT collection and giving them the tokens to buy it
        saleRef.purchase(tokenID: 1, recipient: self.collectionRef, buyTokens: <-self.temporaryVault)

        log("Token 1 has been bought by account 2!")
    }
}
