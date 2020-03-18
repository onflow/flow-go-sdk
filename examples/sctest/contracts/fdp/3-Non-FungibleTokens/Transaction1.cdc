// Transaction1.cdc

import NonFungibleToken from 0x04

// This is a transaction that checks if an NFT exists in storage.
transaction {
    prepare(acct: Account) {
        if acct.storage[NonFungibleToken.NFT] != nil {
            log("The token exists!")
        } else {
            log("No token found!")
        }
    }
}
