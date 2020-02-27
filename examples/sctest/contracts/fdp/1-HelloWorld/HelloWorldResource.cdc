// HelloWorldResource.cdc

// This contract is another example of a Hello World
// contract, but uses a resource, Cadence's special type
// that models ownership.
access(all) contract HelloWorld {

    // Declare a resource that only includes one function
    //
	access(all) resource HelloAsset {
        // A transaction can call this function to 
        // get the Hello World message from the resource
		access(all) fun hello(): String {
			return "Hello World!"
		}
	}

	init() {

        // We can do anything in the init function, including accessing
        // the storage of the account that this contract is deployed to.
        // Here, we are storing the newly created HelloAsset resource
        // in the private account.storage.
		let oldHello <- self.account.storage[HelloAsset] <- create HelloAsset()

        // We have to move the old value out of storage to ensure that
        // it doesn't get accidentally lost of deleted.
        // We can delete it here becuase we know it is empty, but in
        // a real smart contract, we might do something else with it 
        // if it has a value.
		destroy oldHello
	}
}