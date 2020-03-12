// HelloWorld.cdc

// Contracts are deployed by specifying their access level,
// followed by the contract keywork and the name of the contract

// This contract simply holds a string field and has one function
// that returns the contents of the field
access(all) contract HelloWorld {

    // declare a fully public field with type String
    // All fields must be initialized in the init() function
    access(all) let greeting: String

    // The init function is mandatory if there are any fields
    // in the composite type. Here, we simply initialize the 
    // greeting field to "Hello World!"
    init() {
        self.greeting = "Hello World!"
    }

    // Public function that returns our friendly greeting!
    access(all) fun hello(): String {
        return self.greeting
    }
}