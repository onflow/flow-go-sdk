// HelloWorldResource.cdc
//
// This is a variation of the HelloWorld contract that introduces the concept of
// resources, a new form of linear type that is unique to Cadence. Resources can be
// used to create a secure model of digital ownership.
//
// Learn more about resources in this tutorial: https://docs.onflow.org/docs/hello-world

access(all) contract HelloWorld {

    // Declare a resource that only includes one function.
	access(all) resource HelloAsset {
        // A transaction can call this function to get the "Hello, World!"
        // message from the resource.
		access(all) fun hello(): String {
			return "Hello, World!"
		}
	}

	init() {
        // We can do anything in the init function, including accessing
        // the storage of the account that this contract is deployed to.
        //
        // Here we are storing the newly created HelloAsset resource
        // in the private account.storage.
		let oldHello <- self.account.storage[HelloAsset] <- create HelloAsset()

        // We have to move the old value out of storage to ensure that
        // it doesn't get accidentally lost or deleted.
        //
        // We can delete it here because we know that it's empty, but in
        // a real smart contract we might do something else with it
        // if it has a value.
		destroy oldHello

        log("HelloAsset created and stored")
	}
}
