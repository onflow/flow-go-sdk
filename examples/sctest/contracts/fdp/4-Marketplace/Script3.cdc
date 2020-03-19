// Script3.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02
import Marketplace from 0x03

// This script checks to make sure that the accounts'
// Vault balances and NFT collections are correct
//
// Account 1: Vault balance = 50, No NFTs
// Account 2: Vault balance = 10, NFT ID=1
pub fun main() {
    // Get both accounts' public account objects
    let account1 = getAccount(0x01)
	let account2 = getAccount(0x02)

    // Get both accounts' published Vault Receiver references 
    let acct1ftRef = account1.published[&FungibleToken.Receiver] ?? panic("missing account 1 vault reference")
    let acct2ftRef = account2.published[&FungibleToken.Receiver] ?? panic("missing account 2 vault reference")

    // Log both accounts' Vault balances and ensure they are
    // the correct numbers. Account 1 should have 50 and
    // Account 2 should have 10
    log("Account 1 Vault Balance")
    log(acct1ftRef.balance)
    log("Account 2 Vault Balance")
    log(acct2ftRef.balance)
    if acct1ftRef.balance != UInt64(50) || acct2ftRef.balance != UInt64(10) {
        panic("Wrong Balances!")
    }

    // Get both accounts' published NFT Collection Receiver references
    let acct1nftRef = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 1 nft reference!")
	let acct2nftRef = account2.published[&NonFungibleToken.NFTReceiver] ?? panic("missing account 2 nft reference!")

    // Log the NFT ids that both accounts own
    // Account 1 should have none
    // Account 2 should have NFT 1
    log("Account 1 NFTs")
    log(acct1nftRef.getIDs())
	log("Account 2 NFTs")
    log(acct2nftRef.getIDs())

    if acct1nftRef.getIDs().length != 0 || acct2nftRef.getIDs()[0] != UInt64(1) {
        panic("Wrong NFTs in Collection")
    }

    // Get the public sale reference for Account 0x01
    let acct1SaleRef = account1.published[&Marketplace.SalePublic] ?? panic("missing account 1 Sale reference!")

    // print the NFTs that account 1 has for sale
    log("Account 1 NFTs for Sale")
    log(acct1SaleRef.getIDs())
    if acct1SaleRef.getIDs().length != 0 { panic("Sale should be empty!") }

}