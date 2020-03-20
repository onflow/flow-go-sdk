// Script1.cdc

import FungibleToken from 0x01 

// This script reads the Vault balances of two accounts.
access(all) fun main() {
    // Get the accounts' public account objects
    let acct1 = getAccount(0x01)
    let acct2 = getAccount(0x02)

    // Use optional chaining to read and log balance fields
    log("Account 1 Balance")
	log(acct1.published[&FungibleToken.Receiver]?.balance)
    log("Account 2 Balance")
    log(acct2.published[&FungibleToken.Receiver]?.balance)
}
