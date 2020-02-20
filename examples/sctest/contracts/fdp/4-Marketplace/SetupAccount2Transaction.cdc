// SetupAccount2Transaction.cdc

import FungibleToken from 0x01
import NonFungibleToken from 0x02

transaction {
        prepare(acct: Account) {
            // create a new vault instance with an initial balance of 30
            let vaultA <- FungibleToken.createEmptyVault()
            
            // store it in the account storage
            // and destroy whatever was there previously
            let oldVault <- acct.storage[FungibleToken.Vault] <- vaultA
            destroy oldVault
            acct.published[&FungibleToken.Receiver] = &acct.storage[FungibleToken.Vault] as &FungibleToken.Receiver

            // publish a public interface that only exposes ownedNFTs and deposit
            acct.published[&NonFungibleToken.NFTReceiver] = &acct.storage[NonFungibleToken.Collection] as &NonFungibleToken.Collection
            let account1 = getAccount(0x01)
            let acct1nftReceiver = account1.published[&NonFungibleToken.NFTReceiver] ?? panic("no receiver found")
            let minterRef = acct.storage[&NonFungibleToken.NFTMinter] ?? panic("no minter found")
            minterRef.mintNFT(recipient: acct1nftReceiver)
        
        }
        execute {}
}
 