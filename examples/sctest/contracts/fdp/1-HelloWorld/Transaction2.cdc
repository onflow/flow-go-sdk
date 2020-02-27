// Transaction2.cdc

import HelloWorld from 0x01

transaction {

    // In this transaction, we are calling the hello method
    // on the HelloAsset object that is stored in account storage
    prepare(acct: Account) {
        
        // we use optional chaining (?) because the value in storage
        // is an optional
        log(acct.storage[HelloWorld.HelloAsset]?.hello())
    }
}