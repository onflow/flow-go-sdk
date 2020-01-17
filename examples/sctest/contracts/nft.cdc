
pub contract interface NonFungibleToken {

    // The total number of tokens of this type in existance
    pub var totalSupply: Int

    // event ContractInitialized()
    // event Withdraw()
    // event Deposit()

    pub resource interface INFT {
        // The unique ID that each NFT has
        pub let id: Int

        // placeholder for token metadata 
        pub var metadata: {String: String}
    }

    pub resource NFT: INFT {
        pub let id: Int

        pub var metadata: {String: String}

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

        pub fun batchWithdraw(ids: [Int]): @Collection
    }

    pub resource interface Receiver {

		pub fun deposit(token: @NFT) {
			pre {
				token.id > 0:
					"deposit: token ID must be positive!"
			}
		}

        pub fun batchDeposit(tokens: @Collection)
    }

    pub resource interface Metadata {

		pub fun getIDs(): [Int]

		pub fun idExists(id: Int): Bool {
			pre {
				id > 0: "idExists: token id must be positive!"
			}
		}

        pub fun getMetaData(id: Int, field: String): String {
            pre {
				id > 0: "idExists: token id must be positive!"
                field.length != 0: "getMetaData: field is undefined!"
			}
        }
	}

    pub resource Collection: Provider, Receiver, Metadata {
        
        pub var ownedNFTs: @{Int: NFT}

        // withdraw removes an NFT from the collection and moves it to the caller
        pub fun withdraw(tokenID: Int): @NFT 

        pub fun batchWithdraw(ids: [Int]): @Collection

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @NFT)

        pub fun batchDeposit(tokens: @Collection)

        // idExists checks to see if a NFT with the given ID exists in the collection
        pub fun idExists(tokenID: Int): Bool 

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [Int]

        pub fun getMetaData(id: Int, field: String): String
    }
}

pub contract CryptoKitties: NonFungibleToken {

    pub var totalSupply: Int

    pub resource NFT: NonFungibleToken.INFT {
        pub let id: Int

        pub var metadata: {String: String}

        init(initID: Int) {
            self.id = initID
            self.metadata = {}
        }
    }

    pub resource Collection: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.Metadata {
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

        pub fun batchWithdraw(ids: [Int]): @Collection {
            var i = 0
            var batchCollection: @Collection <- create Collection()

            while i < ids.length {
                batchCollection.deposit(token: <-self.withdraw(tokenID: ids[i]))

                i = i + 1
            }
            return <-batchCollection
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @NFT) {
            let id: Int = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            destroy oldToken
        }

        pub fun batchDeposit(tokens: @Collection) {
            var i = 0
            let keys = tokens.getIDs()

            while i < keys.length {
                self.deposit(token: <-tokens.withdraw(tokenID: keys[i]))

                i = i + 1
            }
            destroy tokens
        }

        // idExists checks to see if a NFT with the given ID exists in the collection
        pub fun idExists(tokenID: Int): Bool {
            return self.ownedNFTs[tokenID] != nil
        }

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [Int] {
            return self.ownedNFTs.keys
        }

        pub fun getMetaData(id: Int, field: String): String {
            let token <- self.ownedNFTs[id] ?? panic("No NFT!")
            
            let data = token.metadata[field] ?? panic("No metadata!")

            let oldToken <- self.ownedNFTs[id] <- token
            destroy oldToken

            return data
        }

        destroy() {
            destroy self.ownedNFTs
        }
    }

    pub fun createNFT(id: Int): @NFT {
        return <- create NFT(initID: id)
    }

    pub fun createCollection(): @Collection {
        return <- create Collection()
    }

	pub resource NFTFactory {

		// the ID that is used to mint moments
		pub var idCount: Int

		init() {
			self.idCount = 1
		}

		// mintNFT mints a new NFT with a new ID
		// and deposit it in the recipients colelction using their collection reference
		pub fun mintNFT(recipient: &Collection) {

					// create a new NFT
			var newNFT <- create NFT(initID: self.idCount)
			
					// deposit it in the recipient's account using their reference
			recipient.deposit(token: <-newNFT)

					// change the id so that each ID is unique
			self.idCount = self.idCount + 1
		}
	}

	init() {
        self.totalSupply = 0
        
		let oldCollection <- self.account.storage[Collection] <- create Collection()
		destroy oldCollection

		self.account.storage[&Collection] = &self.account.storage[Collection] as Collection
        self.account.published[&NonFungibleToken.Receiver] = &self.account.storage[Collection] as NonFungibleToken.Receiver

		let oldFactory <- self.account.storage[NFTFactory] <- create NFTFactory()
		destroy oldFactory

		self.account.storage[&NFTFactory] = &self.account.storage[NFTFactory] as NFTFactory
	}
}

