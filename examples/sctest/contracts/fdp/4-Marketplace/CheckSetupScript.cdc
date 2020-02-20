// CheckSetupScript.cdc

import FungibleToken from 0x01

import NonFungibleToken from 0x02

access(all) fun main() {
    let account1 = getAccount(0x01)
	let account2 = getAccount(0x02)

    let acct1ftRef = account1.published[&FungibleToken.Receiver] ?? panic("missing account 1 vault reference")
    let acct2ftRef = account2.published[&FungibleToken.Receiver] ?? panic("missing account 2 vault reference")

    log("Account 1 Vault Balance")
    log(acct1ftRef.balance)
    log("Account 2 Vault Balance")
    log(acct2ftRef.balance)
    if acct1ftRef.balance != UInt64(40) || acct2ftRef.balance != UInt64(20) {
        panic("Wrong Balances!")
    }

    let acct1nftRef = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 1 nft reference!")
	let acct2nftRef = account2.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 2 nft reference!")

    log("Account 1 NFTs")
    log(acct1nftRef.getIDs())
	log("Account 2 NFTs")
    log(acct2nftRef.getIDs())
    if acct1nftRef.getIDs()[0] != UInt64(1) || acct2nftRef.getIDs().length != 0 {
        panic("Wrong Balances!")
    }
}