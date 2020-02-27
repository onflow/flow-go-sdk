// Script1.cdc

import FungibleToken from 0x01 

// Script to read the Vault balances of two accounts.
access(all) fun main() {
    // get the accounts' public account objects
    let acct1 = getAccount(0x01)
    let acct2 = getAccount(0x02)

    // Using optional chaining to read the balance fields
    // and print them
    log("Account 1 Balance")
	log(acct1.published[&FungibleToken.Receiver]?.balance)
    log("Account 2 Balance")
    log(acct2.published[&FungibleToken.Receiver]?.balance)
}