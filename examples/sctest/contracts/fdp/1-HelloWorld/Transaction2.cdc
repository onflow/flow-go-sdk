// Transaction2.cdc

import HelloWorld from 0x01

transaction {

    prepare(acct: Account) {
        log(acct.storage[HelloWorld.HelloAsset]?.hello())
    }
}