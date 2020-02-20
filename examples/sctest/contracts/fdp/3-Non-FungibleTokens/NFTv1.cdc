// NFTv1.cdc

access(all) contract NonFungibleToken {

    access(all) resource NFT {
        access(all) let id: UInt64

        access(all) var metadata: {String: String}

        init(initID: UInt64) {
            self.id = initID
            self.metadata = {}
        }
    }

	  init() {
				let oldNFT <- self.account.storage[NFT] <- create NFT(initID: 1)
				destroy oldNFT
		}
}