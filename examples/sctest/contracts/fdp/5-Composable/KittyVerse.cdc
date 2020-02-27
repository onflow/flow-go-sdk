// KittyVerse.cdc

// KittyVerse is a contract that defines two types of NFTs
// One is a KittyHat, which represents a special Hat
// The second is the Kitty resource, which can own Kitty Hats
//
access(all) contract KittyVerse {

    access(all) resource KittyHat {
        access(all) let id: Int
        access(all) let name: String

        init(id: Int, name: String) {
            self.id = id
            self.name = name
        }

        // An example of a function someone might put in their hat resource
        access(all) fun tipHat(): String {
            if self.name == "Cowboy Hat" {
                return "Howdy Y'all"
            } else if self.name == "Top Hat" {
                return "Greetings, fellow aristocats!"
            } 

            return "Hello"
        }
    }

    access(all) fun createHat(id: Int, name: String): @KittyHat {
        return <-create KittyHat(id: id, name: name)
    }

    access(all) resource Kitty {

        access(all) let id: Int

        // place where the Kitty hats are stored
        access(all) let items: @{String: KittyHat}

        init(newID: Int) {
            self.id = newID
            self.items <- {}
        }

        destroy() {
            destroy self.items
        }
    }

    access(all) fun createKitty(): @Kitty {
        return <-create Kitty(newID: 1)
    }
}