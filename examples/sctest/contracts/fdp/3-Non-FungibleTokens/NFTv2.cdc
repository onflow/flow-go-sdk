// NFTv2.cdc

access(all) contract NonFungibleToken {

    access(all) resource NFT {
        access(all) let id: UInt64

        access(all) var metadata: {String: String}

        init(initID: UInt64) {
            self.id = initID
            self.metadata = {}
        }
    }

    access(all) resource interface NFTReceiver {
        // dictionary of NFT conforming tokens
        access(all) var ownedNFTs: @{UInt64: NFT}

        access(all) fun deposit(token: @NFT)

        access(all) fun getIDs(): [UInt64]

        access(all) fun idExists(id: UInt64): Bool
    }

    access(all) resource Collection: NFTReceiver {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `UInt64` ID field
        access(all) var ownedNFTs: @{UInt64: NFT}

        init () {
            self.ownedNFTs <- {}
        }

        // withdraw removes an NFT from the collection and moves it to the caller
        access(all) fun withdraw(withdrawID: UInt64): @NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            return <-token
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        access(all) fun deposit(token: @NFT) {
            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token
            destroy oldToken
        }

        // idExists checks to see if a NFT with the given ID exists in the collection
        access(all) fun idExists(id: UInt64): Bool {
            return self.ownedNFTs[id] != nil
        }

        // getIDs returns an array of the IDs that are in the collection
        access(all) fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        destroy() {
            destroy self.ownedNFTs
        }
    }

    access(all) fun createEmptyCollection(): @Collection {
        return <- create Collection()
    }

	init() {
		// store an empty NFT Collection in account storage
        let oldCollection <- self.account.storage[Collection] <- create Collection()
        destroy oldCollection
	}
}