// Transaction1.cdc

import HelloWorld from 0x01

transaction {

	  prepare(acct: Account) {}
	
	  execute {
	  	  log(HelloWorld.hello())
	  }
}