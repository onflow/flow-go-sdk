// NFTv3.cdc

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
            let id = token.id

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

    access(all) resource NFTMinter {

        // the ID that is used to mint NFTs
        access(all) var idCount: UInt64

        init() {
            self.idCount = 1
        }

        // mintNFT mints a new NFT with a new ID
        // and deposit it in the recipients colelction using their collection reference
        access(all) fun mintNFT(recipient: &NFTReceiver) {

            // create a new NFT
            var newNFT <- create NFT(initID: self.idCount)
            
            // deposit it in the recipient's account using their reference
            recipient.deposit(token: <-newNFT)

            // change the id so that each ID is unique
            self.idCount = self.idCount + UInt64(1)
        }
    }

	init() {
		// store an empty NFT Collection in account storage
        let oldCollection <- self.account.storage[Collection] <- create Collection()
        destroy oldCollection

        self.account.published[&NFTReceiver] = &self.account.storage[Collection] as &NFTReceiver

        let oldMinter <- self.account.storage[NFTMinter] <- create NFTMinter()
        destroy oldMinter

        self.account.storage[&NFTMinter] = &self.account.storage[NFTMinter] as &NFTMinter
	}
}
 