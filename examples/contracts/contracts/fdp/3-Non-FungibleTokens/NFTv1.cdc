// NFTv1.cdc
//
// The NonFungibleToken contract is a sample implementation of a non-fungible token (NFT) on Flow.
//
// This contract defines one of the simplest forms of NFTs using an
// integer ID and metadata field.
// 
// Learn more about non-fungible tokens in this tutorial: https://docs.onflow.org/docs/non-fungible-tokens

access(all) contract NonFungibleToken {

    // Declare the NFT resource type
    access(all) resource NFT {
        // The unique ID that differentiates each NFT
        access(all) let id: UInt64

        // String mapping to hold metadata
        access(all) var metadata: {String: String}

        // Initialize both fields in the init function
        init(initID: UInt64) {
            self.id = initID
            self.metadata = {}
        }
    }

    // Create a single new NFT and put it into account storage
	init() {
		let oldNFT <- self.account.storage[NFT] <- create NFT(initID: 1)
		destroy oldNFT
	}
}
