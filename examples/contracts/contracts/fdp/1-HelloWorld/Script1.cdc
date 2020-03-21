// Script1.cdc

import HelloWorld from 0x02

// A script is a special type of Cadence transaction
// that does not have access to any account's storage
// and cannot modify state. Any state changes that it would
// do are reverted after execution.
//
// Scripts must have the following signature: access(all) fun main()
access(all) fun main() {

    // Cadence code can get an account's public account object
    // by using the getAccount() built-in function.
	let helloAccount = getAccount(0x02)

    // The log built-in function logs its argument to stdout.
    //
    // Here we are using optional chaining to call the "hello"
    // method on the HelloAsset resource that is referenced
    // in the published area of the account.
	log(helloAccount.published[&HelloWorld.HelloAsset]?.hello())
}
