
pub contract interface INonFungibleToken {

    // The total number of tokens of this type in existance
    pub var totalSupply: Int

    pub resource interface INFT {
        // The unique ID that each NFT has
        pub let id: Int
    }

    pub resource NFT: INFT {
        pub let id: Int

        init(initID: Int) {
            pre {
                initID > 0: "NFT init: initID must be positive!"
            }
        }
    }

    pub resource interface Provider {
        // withdraw removes an NFT from the collection and moves it to the caller
        pub fun withdraw(withdrawID: Int): @NFT {
            pre {
                withdrawID > 0: "withdraw: withdrawID must be positive"
            }
            post {
                result.id == withdrawID: "withdraw: The ID of the withdrawn token must be the same as the requested ID"
            }
        }
    }

    pub resource interface Receiver {

		pub fun deposit(token: @NFT) {
			pre {
				token.id > 0:
					"deposit: token ID must be positive!"
			}
		}
    }

    pub resource interface Metadata {

		pub fun getIDs(): [Int]

		pub fun idExists(id: Int): Bool {
			pre {
				tokenID > 0:
					"idExists: token id must be positive!"
			}
		}
	}
}

pub contract NonFungibleToken {

    pub var totalSupply: Int

    pub resource NFT: INonFungibleToken.INFT {
        pub let id: Int

        init(newID: Int) {
            pre {
                newID > 0: "NFT init: NFT ID must be positive!"
            }
        }
    }

	

    pub resource NFTCollection: Receiver {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `Int` ID field
        pub var ownedNFTs: @{Int: NFT}

        init () {
            self.ownedNFTs <- {}
        }

        // withdraw removes an NFT from the collection and moves it to the caller
        pub fun withdraw(tokenID: Int): @NFT {
            let token <- self.ownedNFTs.remove(key: tokenID) ?? panic("missing NFT")

            return <-token
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @NFT): Void {
            let id: Int = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            destroy oldToken
        }

        // idExists checks to see if a NFT with the given ID exists in the collection
        pub fun idExists(tokenID: Int): Bool {
            return self.ownedNFTs[tokenID] != nil
        }

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [Int] {
            return self.ownedNFTs.keys
        }

        destroy() {
            destroy self.ownedNFTs
        }

        // createCollection returns a new collection resource to the caller
        pub fun createCollection(): @NFTCollection {
            return <- create NFTCollection()
        }
    }

    pub fun createNFT(id: Int): @NFT {
        return <- create NFT(newID: id)
    }

    pub fun createCollection(): @NFTCollection {
        return <- create NFTCollection()
    }

	pub resource NFTFactory {

		// the ID that is used to mint moments
		pub var idCount: Int

		init() {
			self.idCount = 1
		}

		// mintNFT mints a new NFT with a new ID
		// and deposit it in the recipients colelction using their collection reference
		pub fun mintNFT(recipient: &NFTCollection) {

					// create a new NFT
			var newNFT <- create NFT(newID: self.idCount)
			
					// deposit it in the recipient's account using their reference
			recipient.deposit(token: <-newNFT)

					// change the id so that each ID is unique
			self.idCount = self.idCount + 1
		}
	}

	init() {
		let oldCollection <- self.account.storage[NFTCollection] <- create NFTCollection()
		destroy oldCollection

		self.account.storage[&NFTCollection] = &self.account.storage[NFTCollection] as NFTCollection
        self.account.published[&Receiver] = &self.account.storage[NFTCollection] as Receiver

		let oldFactory <- self.account.storage[NFTFactory] <- create NFTFactory()
		destroy oldFactory

		self.account.storage[&NFTFactory] = &self.account.storage[NFTFactory] as NFTFactory
	}
}

