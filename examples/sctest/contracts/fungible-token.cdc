

// The Fungible Token standard interface that all Fungible Tokens
// would have to conform to
pub contract interface IFungibleToken {

    // The total number of tokens in existence
    pub var totalSupply: Int

    // // event that is emmited when the contract is created
    // event TokenContractInitialized(initialSupply: Int, location: Address)
    // // event that is emmited when tokens are withdrawn from a Vault
    // event TokensWithdrawn(amount: Int, recipient: Address)
    // // event that is emitted when tokens are deposited to a Vault
    // event TokensDeposited(amount: Int, recipient: Address)

    // Interface that enforces the requirements for withdrawing
    // tokens from the implementing type
    //
    pub resource interface Provider {
        pub fun withdraw(amount: Int): @Vault {
            pre {
                amount >= 0:
                    "withdraw: Withdrawal amount must be non-negative"
            }
            post {
                result.balance == amount:
                    "withdraw: Incorrect amount withdrawn"
            }
        }
    }

    // Interface that enforces the requirements for depositing
    // tokens into the implementing type
    //
    pub resource interface Receiver {
        pub fun deposit(from: @Vault): Void {
            pre {
                from.balance > 0:
                    "deposit: Deposit balance must be positive"
            }
        }
    }

    // Every Fungible Token contract must define a Vault object that
    // conforms to the Provider and Receiver interfaces
    // and includes these fields and functions
    //
    pub resource Vault: Provider, Receiver {
        // keeps track of the total balance of the accounts tokens
        pub var balance: Int

        init(balance: Int) {
            pre {
                balance >= 0:
                    "init: Initial balance must be non-negative"
            }
            post {
                self.balance == balance:
                    "init: balance must be initialized to the initial balance"
            }
        }

        // withdraw will usually subtract `amount` from the vaults balance and
        // return a vault object with the subtracted balance
        pub fun withdraw(amount: Int): @Vault

        // deposit will usually take a vault object as a parameter and add
        // its balance to the balance of the stored vault, then
        // destroy the sent vault because its balance has been consumed
        pub fun deposit(from: @Vault): Void {
            post {
                self.balance == before(self.balance) + before(from.balance):
                    "deposit: Incorrect amount removed"
            }
        }

        // In order to destroy a Vault, its balance must be zero
        // so tokens aren't lost
        //
        destroy() {
            pre {
                self.balance == 0: "destroy: balance must be zero"
            }
        }
    }

    // Any user can call this function to create a new Vault object
    // that has balance = 0
    //
    pub fun createEmptyVault(): @Vault {
        post {
            result.balance == 0: "createEmptyVault: The newly created Vault must have zero balance"
        }
    }
}

pub contract FlowToken: IFungibleToken {

    pub var totalSupply: Int

    pub resource Vault: IFungibleToken.Provider, IFungibleToken.Receiver {
        
        pub var balance: Int

        init(balance: Int) {
            self.balance = balance
        }

        pub fun withdraw(amount: Int): @Vault {
            self.balance = self.balance - amount
            return <-create Vault(balance: amount)
        }
        
        pub fun deposit(from: @Vault): Void {
            self.balance = self.balance + from.balance
            destroy from
        }
    }

    pub fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0)
    }

    pub fun createVault(initialBalance: Int): @Vault {
        return <-create Vault(balance: initialBalance)
    }

    init() {
        self.totalSupply = 1000

        let oldVault <- self.account.storage[Vault] <- create Vault(balance: 1000)
        destroy oldVault

        self.account.storage[&Vault] = &self.account.storage[Vault] as Vault
        self.account.published[&IFungibleToken.Receiver] = &self.account.storage[Vault] as IFungibleToken.Receiver
    }
}
