// CheckSetupScript.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

// Script that makes sure that the accounts are set up correctly
// to complete the Marketplace tutorialx
access(all) fun main() {
    // Get both accounts' public account objects
    let account1 = getAccount(0x01)
	let account2 = getAccount(0x02)

    // get both accounts' published Vault Receiver references 
    let acct1ftRef = account1.published[&FungibleToken.Receiver] ?? panic("missing account 1 vault reference")
    let acct2ftRef = account2.published[&FungibleToken.Receiver] ?? panic("missing account 2 vault reference")

    // Log both accounts' Vault balances and ensure they are
    // the correct numbers. Account 1 should have 40 and
    // Account 2 should have 20
    log("Account 1 Vault Balance")
    log(acct1ftRef.balance)
    log("Account 2 Vault Balance")
    log(acct2ftRef.balance)
    if acct1ftRef.balance != UInt64(40) || acct2ftRef.balance != UInt64(20) {
        panic("Wrong Balances!")
    }

    // Get both accounts' published NFT Collection Receiver references
    let acct1nftRef = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 1 nft reference!")
	let acct2nftRef = account2.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 2 nft reference!")

    // Log the NFT ids that both accounts own
    // Account 1 should have NFT 1
    // Account 2 should have none
    log("Account 1 NFTs")
    log(acct1nftRef.getIDs())
	log("Account 2 NFTs")
    log(acct2nftRef.getIDs())
    if acct1nftRef.getIDs()[0] != UInt64(1) || acct2nftRef.getIDs().length != 0 {
        panic("Wrong Balances!")
    }
}