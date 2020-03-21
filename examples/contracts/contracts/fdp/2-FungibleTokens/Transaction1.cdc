// Transaction1.cdc

import FungibleToken from 0x01

// This transaction publishes a reference to the account's token vault. The reference
// can only be used to deposit funds into the account.
transaction {
    prepare(acct: AuthAccount) {

        // Cast the Vault as a FungibleToken.Receiver interface, which only exposes the
        // balance field and deposit function of the underlying vault.
        acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

        log("Public Receiver reference created!")
    }

    post {
        // Check that the reference was created correctly
        getAccount(0x01).published[&FungibleToken.Receiver] != nil:  "Vault Receiver Reference was not created correctly"
    }
}
