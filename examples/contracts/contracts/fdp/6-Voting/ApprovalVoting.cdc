/*
* 
*   In this example, we want to create a simple approval voting contract 
*   where a polling place issues ballots to addresses. 
*   
*   The run a vote, the Admin deploys the smart contract,
*   then initializes the proposals 
*   using the initialize_proposals.cdc transaction.
*   The array of proposals cannot be modified after it has been initialized.
*
*   Then they will give ballots to users by 
*   using the issue_ballot.cdc transaction.
*
*   Every user with a ballot is allowed to approve any number of proposals. 
*   A user can choose their votes and cast them 
*   with the cast_vote.cdc transaction.
* 
*/

pub contract ApprovalVoting {

    //list of proposals to be approved
    pub var proposals: [String]

    // number of votes per proposal
    pub let votes: {Int: Int}

    // This is the resource that is issued to users.
    // When a user gets a Ballot object, they call the `vote` function
    // to include their votes, and then cast it in the smart contract 
    // using the `cast` function to have their vote included in the polling
    pub resource Ballot {

        // array of all the proposals 
        pub let proposals: [String]

        // corresponds to an array index in proposals after a vote
        pub var choices: {Int: Bool}

        init() {
            self.proposals = ApprovalVoting.proposals
            self.choices = {}
        }

        // modifies the ballot
        // to indicate which proposals it is voting for
        pub fun vote(proposal: Int) {
            pre {
                self.proposals[proposal] != nil: "Cannot vote for a proposal that doesn't exist"
            }
            self.choices[proposal] = true
        }
    }

    // Resource that the Administrator of the vote controls to
    // initialize the proposals and to pass out ballot resources to voters
    pub resource Administrator {

        // function to initialize all the proposals for the voting
        pub fun initializeProposals(_ proposals: [String]) {
            pre {
                ApprovalVoting.proposals.length == 0: "Proposals can only be initialized once"
                proposals.length > 0: "Cannot initialize with no proposals"
            }
            ApprovalVoting.proposals = proposals
        }

        // The admin calls this function to create a new Ballo
        // that can be transferred to another user
        pub fun issueBallot(): @Ballot {
            return <-create Ballot()
        }
    }

    // A user moves their ballot to this function in the contract where 
    // its votes are tallied and the ballot is destroyed
    pub fun cast(ballot: @Ballot) {
        var index = 0
        // look through the ballot
        while index < self.proposals.length {
            if ballot.choices[index]! {
                // tally the vote if it is approved
                self.votes[index] = self.votes[index]! + 1
            }
            index = index + 1;
        }
        // Destroy the ballot because it has been tallied
        destroy ballot
    }

    // initializes the contract by setting the proposals and votes to empty 
    // and creating a new Admin resource to put in storage
    init() {
        self.proposals = []
        self.votes = {}

        self.account.storage[Administrator] <-! create Administrator()
    }
}
