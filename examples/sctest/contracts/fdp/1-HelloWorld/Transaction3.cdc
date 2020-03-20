// Transaction3.cdc

import HelloWorld from 0x02

transaction {

    // This transaction creates a new reference to the HelloAsset resource in storage
    // and adds it to the account's published area.
	prepare(account: AuthAccount) {

        // If an area in storage is prefixed by the "&" symbol,
        // it means that it is a reference to an object, not the object itself.
        account.published[&HelloWorld.HelloAsset] = &account.storage[HelloWorld.HelloAsset] as &HelloWorld.HelloAsset

        // Call the hello function using the reference to the HelloResource resource.
        //
        // We use the "?" symbol because the value we are accessing is an optional.
        log(account.published[&HelloWorld.HelloAsset]?.hello())
	}
}
