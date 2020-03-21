// Transaction2.cdc

import HelloWorld from 0x02

transaction {

    // This transaction calls the "hello" method on the HelloAsset object
    // that is stored in the account's storage.
    prepare(acct: AuthAccount) {

        // We use optional chaining (?) because the value in storage
        // may or may not exist, and thus is considered optional.
        log(acct.storage[HelloWorld.HelloAsset]?.hello())
    }
}
