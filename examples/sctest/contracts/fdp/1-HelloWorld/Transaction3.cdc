// Transaction3.cdc

import HelloWorld from 0x01

transaction {

	prepare(acct: Account) {
        // Create a new reference to the HelloAsset resource in Storage
        // and store it in published
        // if an area in storage is prefixed by the & symbol
        // it means that it is a reference to an object, not the object itself
        acct.published[&HelloWorld.HelloAsset] = &acct.storage[HelloWorld.HelloAsset] as &HelloWorld.HelloAsset

        // Call the hello function using the reference to the HelloResource resource
        // We use the ? symbol because the value we are accessing is an optional
        log(acct.published[&HelloWorld.HelloAsset]?.hello())
	}
}