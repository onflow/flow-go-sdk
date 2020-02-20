// Script2.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02
import Marketplace from 0x03

pub fun main() {
    let account1 = getAccount(0x01)

    let acct1saleRef = account1.published[&Marketplace.SalePublic] ?? nil

    log("Account 1 NFTs for sale")
    log(acct1saleRef?.getIDs())
    log("Price")
    log(acct1saleRef?.idPrice(tokenID: 1))
}