// NFTv1.cdc

// This contract defines the simplest form of an NFT 
// with an integer ID and metadata field
// 
// Users would transfer it by directly moving it from
// one account's storage to the other
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

    // Create a single new NFT and put it in account storage
	init() {
		let oldNFT <- self.account.storage[NFT] <- create NFT(initID: 1)
		destroy oldNFT
	}
}