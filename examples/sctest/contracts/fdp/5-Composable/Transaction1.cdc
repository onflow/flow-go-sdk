// Transaction1.cdc

import KittyVerse from 0x01

// This transaction creates a new kitty, creates two new hats
// puts the hats on the cat, and then stores the kitty in accouht storage

transaction {

		prepare(acct: Account) {

				// Create the Kitty object
				let kitty <- KittyVerse.createKitty()

				// Create the KittyHat objects
				let hat1 <- KittyVerse.createHat(id: 1, name: "Cowboy Hat")
				let hat2 <- KittyVerse.createHat(id: 2, name: "Top Hat")

				// Put the hat on the cat!
				let oldCowboyHat <- kitty.items["Cowboy Hat"] <- hat1
				destroy oldCowboyHat
				let oldTopHat <- kitty.items["Top Hat"] <- hat2
				destroy oldTopHat

				log("The Cat has the Hats")
		
				// Store the Kitty in storage
				let oldKitty <- acct.storage[KittyVerse.Kitty] <- kitty
				destroy oldKitty
		}
}