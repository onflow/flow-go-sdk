// Transaction2.cdc

import KittyVerse from 0x01

transaction {

		prepare(acct: Account) {

				// move the Kitty out of storage, which moves its hat along with it
				let kittyOpt <- acct.storage[KittyVerse.Kitty] <- nil
				let kitty <- kittyOpt ?? panic("Kitty doesn't exist!")

				// take the cowboy hat off the Kitty
				let cowboyHatOpt <- kitty.items.remove(key: "Cowboy Hat")
				let cowboyHat <- cowboyHatOpt ?? panic("cowboy hat doesn't exist!")

				// Tip the cowboy hat
				log(cowboyHat.tipHat())
				destroy cowboyHat

				// Tip the top hat that is on the Kitty
				log(kitty.items["Top Hat"]?.tipHat())

				// move the Kitty to storage
				// which moves its hat along with it
				let oldKitty <- acct.storage[KittyVerse.Kitty] <- kitty
				destroy oldKitty

                destroy kittyOpt
                destroy cowboyHatOpt
		}
}