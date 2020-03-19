// Script3.cdc

import NonFungibleToken from 0x02

// Print the NFTs owned by accounts 0x01 and 0x02.
pub fun main() {
    // Get both public account objects
    let account1 = getAccount(0x01)
	let account2 = getAccount(0x02)

    // Get both Receiver references to their Collections
    let acct1Ref = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 1 reference!")
	let acct2Ref = account2.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 2 reference!")

    // Print both collections as arrays of IDs
    log("Account 1 NFTs")
    log(acct1Ref.getIDs())

	log("Account 2 NFTs")
    log(acct2Ref.getIDs())
}
