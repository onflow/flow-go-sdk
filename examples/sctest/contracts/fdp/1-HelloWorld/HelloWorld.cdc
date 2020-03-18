// HelloWorld.cdc
//
// Welcome to Cadence! This is one of the simplest programs you can
// deploy on Flow.
//
// This contract contains a single string field along with a public getter function.
//
// Follow the "Hello, World!" tutorial to learn more: https://docs.onflow.org/docs/hello-world

access(all) contract HelloWorld {

    // Declare a fully public field of type String.
    // All fields must be initialized in the init() function.
    access(all) let greeting: String

    // The init function is mandatory if there are any fields
    // in the contract. Here, we simply initialize the
    // greeting field to "Hello, World!"
    init() {
        self.greeting = "Hello, World!"
    }

    // Public function that returns our friendly greeting!
    access(all) fun hello(): String {
        return self.greeting
    }
}