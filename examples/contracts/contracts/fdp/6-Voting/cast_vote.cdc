import ApprovalVoting from 0x01

// This transaction allows a voter to select the votes they would like to make
// and cast that vote by using the castVote function 
// of the ApprovalVoting smart contract

transaction {
    prepare(voter: AuthAccount) {
        
        // take the voter's ballot our of storage
        let ballot <- voter.storage[ApprovalVoting.Ballot]!

        // Vote on the proposal 
        ballot.vote(proposal: 1)

        // Cast the vote by submitting it to the smart contract
        ApprovalVoting.cast(ballot: <-ballot)
    }
}