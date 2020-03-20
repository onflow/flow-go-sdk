// NFTv2.cdc
//
// This is an extended version of the NonFungibleToken contract
// that includes withdraw and deposit functionality, as well as a
// collection resource that can be used to bundle NFTs together.
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

    // We define this interface purely as a way to allow users
    // to create public, restricted references to their NFT Collection.
    // They would use this to only expose the deposit, getIDs,
    // and idExists fields in their Collection
    access(all) resource interface NFTReceiver {

        access(all) fun deposit(token: @NFT)

        access(all) fun getIDs(): [UInt64]

        access(all) fun idExists(id: UInt64): Bool
    }

    // The definition of the Collection resource that
    // holds the NFTs that a user owns
    access(all) resource Collection: NFTReceiver {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `UInt64` ID field
        access(all) var ownedNFTs: @{UInt64: NFT}

        // Initialize the NFTs field to an empty collection
        init () {
            self.ownedNFTs <- {}
        }

        // withdraw 
        //
        // Function that removes an NFT from the collection 
        // and moves it to the calling context
        access(all) fun withdraw(withdrawID: UInt64): @NFT {
            // If the NFT isn't found, the transaction panics and reverts
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            return <-token
        }

        // deposit 
        //
        // Function that takes a NFT as an argument and 
        // adds it to the collections dictionary
        access(all) fun deposit(token: @NFT) {
            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[token.id] <- token
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

    // creates a new empty Collection resource and returns it 
    access(all) fun createEmptyCollection(): @Collection {
        return <- create Collection()
    }

	init() {
		// store an empty NFT Collection in account storage
        let oldCollection <- self.account.storage[Collection] <- create Collection()
        destroy oldCollection
	}
}
