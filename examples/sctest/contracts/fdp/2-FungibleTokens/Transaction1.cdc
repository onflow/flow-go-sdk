// Transaction1.cdc

import FungibleToken from 0x01
 
transaction {
    prepare(acct: Account) {
        // store a reference to the vault cast as Receiver
        // that is used by external accounts
        // This reference only exposes the balance field and the deposit function
        acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

        log("Public Receiver reference created!")
    }

    post {
        // make sure the reference was created correctly
        getAccount(0x01).published[&FungibleToken.Receiver] != nil:  "Vault Receiver Reference was not created correctly"
    }
}