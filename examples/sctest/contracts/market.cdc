import FungibleToken, FlowToken from 0x0000000000000000000000000000000000000002
import NonFungibleToken, Tokens from 0x0000000000000000000000000000000000000003

// Marketplace is where users can put their NFTs up for sale with a price
// if another user sees an NFT that they want to buy,
// they can send fungible tokens that equal or exceed the buy price
// to buy the NFT.  The NFT is transferred to them when
// they make the purchase


// each user who wants to sell tokens will have a sale collection 
// instance in their account that holds the tokens that they are putting up for sale

// They will give a reference to this collection to the central contract
// that it can use to list tokens

// this is only for one type of NFT for now
// TODO: make it a marketplace that can buy and sell different classes of NFTs
//       how do I do this? Generic NFTs? Using interfaces instead of resources as argument
//       and storage types?


access(all) contract Market {

    access(all) event ForSale(id: UInt64, price: UInt64)
    access(all) event PriceChanged(id: UInt64, newPrice: UInt64)
    access(all) event TokenPurchased(id: UInt64, price: UInt64)
    access(all) event SaleWithdrawn(id: UInt64)

    access(all) resource interface SalePublic {

        access(all) fun purchase(tokenID: UInt64, recipient: &AnyResource{NonFungibleToken.Receiver}, buyTokens: @FungibleToken.Vault)

        access(all) fun idPrice(tokenID: UInt64): UInt64?

        access(all) fun getIDs(): [UInt64]
    }

    access(all) resource SaleCollection: SalePublic {

        // a dictionary of the NFTs that the user is putting up for sale
        access(all) var forSale: @{UInt64: Tokens.NFT}

        // dictionary of the prices for each NFT by ID
        access(all) var prices: {UInt64: UInt64}

        // the fungible token vault of the owner of this sale
        // so that when someone buys a token, this resource can deposit
        // tokens in their account
        access(account) let ownerVault: &AnyResource{FungibleToken.Receiver}

        init (vault: &AnyResource{FungibleToken.Receiver}) {
            self.forSale <- {}
            self.ownerVault = vault
            self.prices = {}
        }

        // withdraw gives the owner the opportunity to remove a sale from the collection
        access(all) fun withdraw(tokenID: UInt64): @Tokens.NFT {
            // remove the price
            self.prices.remove(key: tokenID)
            // remove and return the token
            let token <- self.forSale.remove(key: tokenID) ?? panic("missing NFT")
            return <-token
        }

        // listForSale lists an NFT for sale in this collection
        access(all) fun listForSale(token: @Tokens.NFT, price: UInt64) {
            let id: UInt64 = token.id

            self.prices[id] = price

            let oldToken <- self.forSale[id] <- token

            emit ForSale(id: id, price: price)

            destroy oldToken
        }

        // changePrice changes the price of a token that is currently for sale
        access(all) fun changePrice(tokenID: UInt64, newPrice: UInt64) {
            self.prices[tokenID] = newPrice

            emit PriceChanged(id: tokenID, newPrice: newPrice)
        }

        // purchase lets a user send tokens to purchase an NFT that is for sale
        access(all) fun purchase(tokenID: UInt64, recipient: &AnyResource{NonFungibleToken.Receiver}, buyTokens: @FungibleToken.Vault) {
            pre {
                self.forSale[tokenID] != nil && self.prices[tokenID] != nil:
                    "No token matching this ID for sale!"
                buyTokens.balance >= (self.prices[tokenID] ?? UInt64(0)):
                    "Not enough tokens to by the NFT!"
            }

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
            let price = self.prices[tokenID]
            return price
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
    access(all) fun createSaleCollection(ownerVault: &AnyResource{FungibleToken.Receiver}): @SaleCollection {
        return <- create SaleCollection(vault: ownerVault)
    }
}
