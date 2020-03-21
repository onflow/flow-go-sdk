// Transaction1.cdc

import HelloWorld from 0x01

// Transactions are declared with the "transaction" keyword.
//
// They have three possible stages:
//
// prepare: This block has full access to the accounts that signed
// the transaction. This is where resources should be stored, removed,
// or used to create references.
//
// execute: This is the main part of the transaction and is
// the stage where function calls and calls to external
// contracts and resources should occur.
//
// post: This block contains optional postcondition checks that can be used to
// verify that the transaction was executed correctly.

transaction {

    // No need to do anything in prepare because we are not working with
    // account storage.
	prepare(acct: AuthAccount) {}
	
    // In execute, we simply call the hello function
    // of the HelloWorld contract and log the returned String.
	execute {
	  	log(HelloWorld.hello())
	}
}
