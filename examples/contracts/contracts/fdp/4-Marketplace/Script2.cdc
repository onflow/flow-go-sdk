// Script2.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02
import Marketplace from 0x03

// This script prints the NFTs that account 0x01 has for sale.
pub fun main() {
    // Get the public account object for account 0x01
    let account1 = getAccount(0x01)

    // Find the public Sale reference to their Collection
    let acct1saleRef = account1.published[&Marketplace.SalePublic] ?? nil

    // Los the NFTs that are for sale
    log("Account 1 NFTs for sale")
    log(acct1saleRef?.getIDs())
    log("Price")
    log(acct1saleRef?.idPrice(tokenID: 1))
}
