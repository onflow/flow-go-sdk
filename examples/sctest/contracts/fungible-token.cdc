

// The Fungible Token standard interface that all Fungible Tokens
// would have to conform to
pub contract interface FungibleToken {

    // The total number of tokens in existence
    pub var totalSupply: UInt64

    // event that is emitted when the contract is created
    event FungibleTokenInitialized(initialSupply: UInt64)

    // event that is emitted when tokens are withdrawn from a Vault
    event Withdraw(amount: UInt64)

    // event that is emitted when tokens are deposited to a Vault
    event Deposit(amount: UInt64)

    // Interface that enforces the requirements for withdrawing
    // tokens from the implementing type
    //
    pub resource interface Provider {
        pub fun withdraw(amount: UInt64): @Vault {
            post {
                result.balance == amount:
                    "Withdrawal amount must be the same as the balance of the withdrawn Vault"
            }
        }
    }

    // Interface that enforces the requirements for depositing
    // tokens into the implementing type
    //
    pub resource interface Receiver {
        pub fun deposit(from: @Vault) {
            pre {
                from.balance > 0:
                    "Deposit balance must be positive"
            }
        }
    }

    // Interface that contains the balance field of the Vault
    //
    pub resource interface Balance {
        pub var balance: UInt64
    }

    // Every Fungible Token contract must define a Vault object that
    // conforms to the Provider and Receiver interfaces
    // and includes these fields and functions
    //
    pub resource Vault: Provider, Receiver, Balance {
        // keeps track of the total balance of the accounts tokens
        pub var balance: UInt64

        init(balance: UInt64) {
            post {
                self.balance == balance:
                    "Balance must be initialized to the initial balance"
            }
        }

        // withdraw subtracts `amount` from the vaults balance and
        // returns a vault object with the subtracted balance
        pub fun withdraw(amount: UInt64): @Vault {
            pre {
                self.balance >= amount: "Amount withdrawn must be less than the balance of the Vault!"
            }
        }

        // deposit takes a vault object as a parameter and adds
        // its balance to the balance of the stored vault, then
        // destroys the sent vault because its balance has been consumed
        pub fun deposit(from: @Vault) {
            post {
                self.balance == before(self.balance) + before(from.balance):
                    "New Vault balance must be the sum of the previous balance and the deposited Vault"
            }
        }
    }

    // Any user can call this function to create a new Vault object
    // that has balance = 0
    //
    pub fun createEmptyVault(): @Vault {
        post {
            result.balance == 0: "The newly created Vault must have zero balance"
        }
    }
}


// This is an Example Implementation of the Fungible Token Standard
//
pub contract FlowToken: FungibleToken {

    pub var totalSupply: UInt64

    pub resource Vault: FungibleToken.Provider, FungibleToken.Receiver, FungibleToken.Balance {
        
        pub var balance: UInt64

        init(balance: UInt64) {
            self.balance = balance
        }

        pub fun withdraw(amount: UInt64): @Vault {
            self.balance = self.balance - amount
            return <-create Vault(balance: amount)
        }
        
        pub fun deposit(from: @Vault) {
            self.balance = self.balance + from.balance
            destroy from
        }
    }

    pub fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0)
    }

    pub fun createVault(initialBalance: UInt64): @Vault {
        return <-create Vault(balance: initialBalance)
    }

    init() {
        self.totalSupply = 1000

        let oldVault <- self.account.storage[Vault] <- create Vault(balance: 1000)
        destroy oldVault

        self.account.storage[&Vault] = &self.account.storage[Vault] as Vault
        self.account.published[&FungibleToken.Receiver] = &self.account.storage[Vault] as FungibleToken.Receiver
    }
}
