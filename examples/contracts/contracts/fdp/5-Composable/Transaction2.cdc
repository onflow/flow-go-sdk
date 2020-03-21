// Transaction2.cdc

import KittyVerse from 0x01

// This transaction moves a kitty out of storage, takes the cowboy hat off of the kitty,
// calls its tip hat function, and then moves it back into storage.
transaction {
    prepare(acct: AuthAccount) {

        // Move the Kitty out of storage, which also moves its hat along with it
        let kittyOpt <- acct.storage[KittyVerse.Kitty] <- nil
        let kitty <- kittyOpt ?? panic("Kitty doesn't exist!")

        // Take the cowboy hat off the Kitty
        let cowboyHatOpt <- kitty.items.remove(key: "Cowboy Hat")
        let cowboyHat <- cowboyHatOpt ?? panic("cowboy hat doesn't exist!")

        // Tip the cowboy hat
        log(cowboyHat.tipHat())
        destroy cowboyHat

        // Tip the top hat that is on the Kitty
        log(kitty.items["Top Hat"]?.tipHat())

        // Move the Kitty to storage, which
        // also moves its hat along with it.
        let oldKitty <- acct.storage[KittyVerse.Kitty] <- kitty
        destroy oldKitty

        destroy kittyOpt
        destroy cowboyHatOpt
    }
}
