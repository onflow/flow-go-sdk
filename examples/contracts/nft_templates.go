package contracts

import (
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
)

// GenerateCreateNFTScript Creates a script that instantiates a new
// NFT instance, then creates an NFT collection instance, stores the
// NFT in the collection, stores the collection in memory, then stores a
// reference to the collection. It also makes sure that the token exists
// in the collection after it has been added to.
// The id must be greater than zero
func GenerateCreateNFTScript(tokenAddr flow.Address, id int) []byte {
	template := `
		import NonFungibleToken, Tokens from 0x%s

		transaction {
		  prepare(acct: Account) {
			let tokenA <- Tokens.createNFT(id: UInt64(%d))

			let collection <- Tokens.createEmptyCollection()

			collection.deposit(token: <-tokenA)

			if collection.idExists(id: %d) == false {
				panic("Token ID doesn't exist!")
			}
			
			let oldCollection <- acct.storage[Tokens.Collection] <- collection
			destroy oldCollection

			acct.published[&NonFungibleToken.Collection] = &acct.storage[Tokens.Collection] as &NonFungibleToken.Collection
		  }
		}
	`

	return []byte(fmt.Sprintf(template, tokenAddr, id, id))
}

// GenerateDepositScript creates a script that withdraws an NFT token
// from a collection and deposits it to another collection
func GenerateDepositScript(tokenCodeAddr flow.Address, receiverAddr flow.Address, transferNFTID int) []byte {
	template := `
		import NonFungibleToken, Tokens from 0x%s

		transaction {
		  prepare(acct: Account) {
			let recipient = getAccount(0x%s)

			let collectionRef = acct.published[&NonFungibleToken.Collection] ?? panic("missing NFT collection reference")
			let depositRef = recipient.published[&NonFungibleToken.Collection] ?? panic("missing deposit reference")

			let nft <- collectionRef.withdraw(withdrawID: %d)

			depositRef.deposit(token: <-nft)
		  }
		}
	`

	return []byte(fmt.Sprintf(template, tokenCodeAddr.String(), receiverAddr.String(), transferNFTID))
}

// GenerateInspectCollectionScript creates a script that retrieves an NFT collection
// from storage and makes assertions about an NFT ID that it contains with the idExists
// function, which uses an array of IDs
func GenerateInspectCollectionScript(nftCodeAddr, userAddr flow.Address, nftID int, shouldExist bool) []byte {
	template := `
		import NonFungibleToken, Tokens from 0x%s

		pub fun main() {
		  let acct = getAccount(0x%s)
		  let collectionRef = acct.published[&NonFungibleToken.Collection] ?? panic("missing collection reference")
		
		  if %v {
		    if collectionRef.ownedNFTs[UInt64(%d)] == nil {
			  panic("Token ID doesn't exist!")
			}
		  } else {
			  if collectionRef.ownedNFTs[UInt64(%d)] != nil {
				panic("Token ID shouldn't exist!")
			  }
		  }
		}
	`

	return []byte(fmt.Sprintf(template, nftCodeAddr, userAddr, shouldExist, nftID, nftID))
}

// GenerateInspectKeysScript creates a script that retrieves an NFT collection
// from storage and reads the array of keys in the dictionary
// arrays can't be compared for equality right now so the first two elements are compared
func GenerateInspectKeysScript(nftCodeAddr, userAddr flow.Address, id1, id2 int) []byte {
	template := `
		import NonFungibleToken, Tokens from 0x%s

		pub fun main() {
		  let acct = getAccount(0x%s)
		  let collectionRef = acct.published[&NonFungibleToken.Collection] ?? panic("missing collection reference")
		
		  let array = collectionRef.getIDs()

		  if array[0] != UInt64(%d) || array[1] != UInt64(%d) {
			panic("Keys array is incorrect!")
		  }
		}
	`

	return []byte(fmt.Sprintf(template, nftCodeAddr, userAddr, id1, id2))
}
