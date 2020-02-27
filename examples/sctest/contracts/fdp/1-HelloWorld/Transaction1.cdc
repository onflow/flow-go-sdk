// Transaction1.cdc

import HelloWorld from 0x01

// Transactions are delcared with the transaction keyword
// They have three possible stages:

// Prepare: Has access to the private account storage objects
// of the accounts that signed the transaction. This is where
// resources should be stored, removed, and used to create references.
//
// execute: This is usually the main part of the transaction and is
// the stage where function calls and calls to external
// contracts and resources should happen
//
// post: No logic is allowed in post, but it is where important
// checks can happen to make sure the transaction was executed correctly.
transaction {

    // No need to do anything in prepare because we are not working with
    // account storage
	prepare(acct: Account) {}
	
    // In execute, we simply call the hello function 
    // of the HelloWorld contract and print the returned Strig.
	execute {
	  	log(HelloWorld.hello())
	}
}