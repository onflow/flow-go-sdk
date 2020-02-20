// Transaction3.cdc

import FungibleToken from 0x01

transaction {

		prepare(acct: Account) {
                // get account 0x02's public account object
				let recipient = getAccount(0x02)
	
                // find their published Reciever reference
                // and record their balance
				let receiverRef = recipient.published[&FungibleToken.Receiver] ?? panic("No receiver reference!")
                let beforeBalance = receiverRef.balance

                // Mint 30 new tokens for account 0x02
				acct.storage[FungibleToken.VaultMinter]?.mintTokens(amount: 30, recipient: receiverRef)

                // Make sure their account balance has been
                // increased by 30 
                if receiverRef.balance == beforeBalance + UInt64(30) {
                    log("30 Tokens minted and deposited to account 0x02")
                }
		}
}
