// HelloWorldResource.cdc

access(all) contract HelloWorld {

	  access(all) resource HelloAsset {
		    access(all) fun hello(): String {
			      return "Hello World!"
		    }
	  }

	  init() {
		    let oldHello <- self.account.storage[HelloAsset] <- create HelloAsset()
		    destroy oldHello
	  }
}