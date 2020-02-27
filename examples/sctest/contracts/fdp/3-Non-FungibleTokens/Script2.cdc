// Script2.cdc

import NonFungibleToken from 0x02

// Print the NFTs that account's 1 and 2 own
pub fun main() {
    // get both public account objects
    let account1 = getAccount(0x01)
	let account2 = getAccount(0x02)

    // get both Receiver references to their Collections
    let acct1Ref = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 1 reference!")
	let acct2Ref = account2.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 2 reference!")

    // print both collections as arrays of IDs
    log("Account 1 NFTs")
    log(acct1Ref.getIDs())
	log("Account 2 NFTs")
    log(acct2Ref.getIDs())
}