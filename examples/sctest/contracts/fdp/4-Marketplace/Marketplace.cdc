import FungibleToken from 0x01
import NonFungibleToken from 0x02

// Marketplace.cdc
//
// The Marketplace contract is a sample implementation of an NFT Marketplace on Flow.
//
// This contract allows users to put their NFTs up for sale. Other users
// can purchase these NFTs with fungible tokens.
//
// Learn more about marketplaces in this tutorial: https://docs.onflow.org/docs/composable-smart-contracts-marketplace

access(all) contract Marketplace {

    // Event that is emitted when a new NFT is put up for sale
    access(all) event ForSale(id: UInt64, price: UInt64)

    // Event that is emitted when the price of an NFT changes
    access(all) event PriceChanged(id: UInt64, newPrice: UInt64)
    
    // Event that is emitted when a token is purchased
    access(all) event TokenPurchased(id: UInt64, price: UInt64)

    // Event that is emitted when a seller withdraws their NFT from the sale
    access(all) event SaleWithdrawn(id: UInt64)

    // Interface that users will publish for their Sale collection
    // that only exposes the methods that are supposed to be public
    //
    access(all) resource interface SalePublic {
        access(all) fun purchase(tokenID: UInt64, recipient: &NonFungibleToken.NFTReceiver, buyTokens: @FungibleToken.Vault)
        access(all) fun idPrice(tokenID: UInt64): UInt64?
        access(all) fun getIDs(): [UInt64]
    }

    // SaleCollection
    //
    // NFT Collection object that allows a user to put their NFT up for sale
    // where others can send fungible tokens to purchase it
    //
    access(all) resource SaleCollection: SalePublic {

        // Dictionary of the NFTs that the user is putting up for sale
        access(all) var forSale: @{UInt64: NonFungibleToken.NFT}

        // Dictionary of the prices for each NFT by ID
        access(all) var prices: {UInt64: UInt64}

        // The fungible token vault of the owner of this sale.
        // When someone buys a token, this resource can deposit
        // tokens into their account.
        access(account) let ownerVault: &FungibleToken.Receiver

        init (vault: &FungibleToken.Receiver) {
            self.forSale <- {}
            self.ownerVault = vault
            self.prices = {}
        }

        // withdraw gives the owner the opportunity to remove a sale from the collection
        access(all) fun withdraw(tokenID: UInt64): @NonFungibleToken.NFT {
            // remove the price
            self.prices.remove(key: tokenID)
            // remove and return the token
            let token <- self.forSale.remove(key: tokenID) ?? panic("missing NFT")
            return <-token
        }

        // listForSale lists an NFT for sale in this collection
        access(all) fun listForSale(token: @NonFungibleToken.NFT, price: UInt64) {
            let id = token.id

            // store the price in the price array
            self.prices[id] = price

            // put the NFT into the the forSale dictionary
            let oldToken <- self.forSale[id] <- token
            destroy oldToken

            emit ForSale(id: id, price: price)
        }

        // changePrice changes the price of a token that is currently for sale
        access(all) fun changePrice(tokenID: UInt64, newPrice: UInt64) {
            self.prices[tokenID] = newPrice

            emit PriceChanged(id: tokenID, newPrice: newPrice)
        }

        // purchase lets a user send tokens to purchase an NFT that is for sale
        access(all) fun purchase(tokenID: UInt64, recipient: &NonFungibleToken.NFTReceiver, buyTokens: @FungibleToken.Vault) {
            pre {
                self.forSale[tokenID] != nil && self.prices[tokenID] != nil:
                    "No token matching this ID for sale!"
                buyTokens.balance >= (self.prices[tokenID] ?? UInt64(0)):
                    "Not enough tokens to by the NFT!"
            }

            // get the value out of the optional
            if let price = self.prices[tokenID] {
                self.prices[tokenID] = nil
                
                // deposit the purchasing tokens into the owners vault
                self.ownerVault.deposit(from: <-buyTokens)

                // deposit the NFT into the buyers collection
                recipient.deposit(token: <-self.withdraw(tokenID: tokenID))

                emit TokenPurchased(id: tokenID, price: price)
            }
        }

        // idPrice returns the price of a specific token in the sale
        access(all) fun idPrice(tokenID: UInt64): UInt64? {
            return self.prices[tokenID]
        }

        // getIDs returns an array of token IDs that are for sale
        access(all) fun getIDs(): [UInt64] {
            return self.forSale.keys
        }

        destroy() {
            destroy self.forSale
        }
    }

    // createCollection returns a new collection resource to the caller
    access(all) fun createSaleCollection(ownerVault: &FungibleToken.Receiver): @SaleCollection {
        return <- create SaleCollection(vault: ownerVault)
    }
}
