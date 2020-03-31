import ApprovalVoting from 0x01

// This transaction allows the administrator of the Voting contract
// to create a new ballot and store it in a voter's account
// The voter and the administrator have to both sign the transaction
// so it can access their storage

transaction {
    prepare(admin: AuthAccount, voter: AuthAccount) {
        
        // create a new Ballot by using the stored administrator resource
        let ballot <- admin.storage[ApprovalVoting.Administrator]?.issueBallot()

        // store that ballot in the voter's account storage
        voter.storage[ApprovalVoting.Ballot] <-! ballot
    }
}