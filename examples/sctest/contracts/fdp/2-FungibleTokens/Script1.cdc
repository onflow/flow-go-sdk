// Script1.cdc

import FungibleToken from 0x01 

// Read account balances
access(all) fun main() {
    let acct1 = getAccount(0x01)
    let acct2 = getAccount(0x02)

    log("Account 1 Balance")
	log(acct1.published[&FungibleToken.Receiver]?.balance)
    log("Account 2 Balance")
    log(acct2.published[&FungibleToken.Receiver]?.balance)
}