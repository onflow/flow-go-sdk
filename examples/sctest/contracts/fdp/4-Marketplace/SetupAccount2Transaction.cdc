// SetupAccount2Transaction.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// This transaction adds an empty Vault to account 0x02
// and mints an NFT with id=1 that is deposited into
// the NFT collection on account 0x01.
transaction {

    // public receiver for account 1's NFT Collection
    let acct1nftReceiver: &NonFungibleToken.NFTReceiver

    // Private reference to this account's minter resource
    let minterRef: &NonFungibleToken.NFTMinter

    prepare(acct: AuthAccount) {
        // create a new vault instance with an initial balance of 30
        let vaultA <- FungibleToken.createEmptyVault()

        // store it in the account storage
        // and destroy whatever was there previously
        let oldVault <- acct.storage[FungibleToken.Vault] <- vaultA
        destroy oldVault
        // publish a receiver reference to the stored Vault
        acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

        log("Created a Vault and published a reference")

        // publish a public interface that only exposes ownedNFTs and deposit
        acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.Collection

        // get account 1's public account object
        let account1 = getAccount(0x01)

        // Get the NFT receiver public reference from account 1
        self.acct1nftReceiver = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("no receiver found")

        // Get the Minter reference from account storage for account 2
        self.minterRef = &acct.storage[NonFungibleToken.NFTMinter] as &NonFungibleToken.NFTMinter
    }
    execute {

        // Mint an NFT and deposit it into account 0x01's collection
        self.minterRef.mintNFT(recipient: self.acct1nftReceiver)

        log("New NFT minted for account 1")
    }
}
