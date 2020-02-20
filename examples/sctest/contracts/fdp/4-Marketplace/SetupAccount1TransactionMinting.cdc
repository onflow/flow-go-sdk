// SetupAccount1TransactionMinting.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

transaction {
        prepare(acct: Account) {
        let account2 = getAccount(0x02)
        let acct1Ref = acct.published[&FungibleToken.Receiver] ?? panic("no receiver 1 Ref")
        let acct2Ref = account2.published[&FungibleToken.Receiver] ?? panic("no receiver 2 Ref")
        let minterRef = &acct.storage[FungibleToken.VaultMinter] as &FungibleToken.VaultMinter
        minterRef.mintTokens(amount: 20, recipient: acct2Ref)
        minterRef.mintTokens(amount: 10, recipient: acct1Ref)
    }
        execute {}
}
