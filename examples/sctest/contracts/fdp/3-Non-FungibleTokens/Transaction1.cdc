// Transaction1.cdc

import NonFungibleToken from 0x04

// Transaction that simply checks to see if an NFT exists in storage
transaction {
    prepare(acct: Account) {
        if acct.storage[NonFungibleToken.NFT] != nil {
            log("The token exists!")
        } else {
            log("No token found!")
        }
    }
}