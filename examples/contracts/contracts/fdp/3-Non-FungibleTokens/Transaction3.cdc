// Transaction3.cdc

import NonFungibleToken from 0x02

// This transaction allows the Minter account to mint an NFT
// and deposit it into its collection.
transaction {

    // The reference to the collection that will be receiving the NFT
    let receiverRef: &NonFungibleToken.NFTReceiver

    // The reference to the Minter resource stored in account storage
    let minterRef: &NonFungibleToken.NFTMinter

    prepare(acct: AuthAccount) {
        // Get the owner's collection reference
        self.receiverRef = acct.published[&NonFungibleToken.NFTReceiver] ?? panic("No receiver")
        
        // Create a Reference to the minter resource
        self.minterRef = &acct.storage[NonFungibleToken.NFTMinter] as &NonFungibleToken.NFTMinter
    }

    execute {
        // Use the minter reference to mint an NFT, which deposits
        // the NFT into the collection that is sent as a parameter.
        self.minterRef.mintNFT(recipient: self.receiverRef)

        log("NFT Minted and deposited to Account 2's Collection")
    }
}
