// Script2.cdc

import NonFungibleToken from 0x02

// Print Account NFTs
pub fun main() {
    let account1 = getAccount(0x01)
	let account2 = getAccount(0x02)

    let acct1Ref = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 1 reference!")
	let acct2Ref = account2.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 2 reference!")

    log("Account 1 NFTs")
    log(acct1Ref.getIDs())
	log("Account 2 NFTs")
    log(acct2Ref.getIDs())
}