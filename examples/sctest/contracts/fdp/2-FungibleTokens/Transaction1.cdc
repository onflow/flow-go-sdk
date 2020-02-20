// Transaction1.cdc

import FungibleToken from 0x01
 
transaction {
    prepare(acct: Account) {
        // store a reference to the vault
        // that is used by external accounts (Receiver)
        acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver
        if acct.published[&FungibleToken.Receiver] != nil {
            log("Vault Receiver Reference created")
        }
    }
}