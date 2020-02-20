// Transaction4.cdc

import FungibleToken from 0x01

transaction {

    prepare(acct: Account) {
        // create a reference to the stored Vault
        let vault = &acct.storage[FungibleToken.Vault] as &FungibleToken.Vault
        
        // withdraw tokens from the signer's account
        let tokens <- vault.withdraw(amount: 10)
    
        // get the other account's public account object
        let recipient = getAccount(0x01)
        
        // fetch the recipient's published receiver Vault reference
        let receiverRef = recipient.published[&FungibleToken.Receiver] ?? panic("missing Vault receiver reference")

        // use the recipients reference to deposit the tokens into their account
        receiverRef.deposit(from: <-tokens)
    }

}
