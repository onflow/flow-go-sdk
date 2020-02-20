// Script1.cdc

import HelloWorld from 0x01

access(all) fun main() {
	let helloAccount = getAccount(0x01)

	log(helloAccount.published[&HelloWorld.HelloAsset]?.hello())
}