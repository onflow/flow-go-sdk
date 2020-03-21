// Transaction1.cdc

import NonFungibleToken from 0x04

// This transaction checks if an NFT exists in the storage of the given account.
transaction {
    prepare(acct: AuthAccount) {
        if acct.storage[NonFungibleToken.NFT] != nil {
            log("The token exists!")
        } else {
            log("No token found!")
        }
    }
}
