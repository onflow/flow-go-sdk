// SetupAccount2Transaction.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// This transaction sets up a Vault for account 2
// and mints a token that gets deposited into
// account 1's collection
transaction {

        // public receiver for account 1's NFT Collection
        let acct1nftReceiver: &NonFungibleToken.NFTReceiver

        // Private reference to this account's minter resource
        let minterRef: &NonFungibleToken.NFTMinter

        prepare(acct: Account) {
            // create a new vault instance with an initial balance of 30
            let vaultA <- FungibleToken.createEmptyVault()
            
            // store it in the account storage
            // and destroy whatever was there previously
            let oldVault <- acct.storage[FungibleToken.Vault] <- vaultA
            destroy oldVault
            // publish a receiver reference to the stored Vault
            acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

            // publish a public interface that only exposes ownedNFTs and deposit
            acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.Collection
            
            // get account 1's public account object
            let account1 = getAccount(0x01)
            self.acct1nftReceiver = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("no receiver found")
            self.minterRef = &acct.storage[NonFungibleToken.NFTMinter] as &NonFungibleToken.NFTMinter
        }
        execute {
            self.minterRef.mintNFT(recipient: self.acct1nftReceiver)
        }
}
 