package consensus
type pbft interface {
	StartConsensus(request *RequestMsg) (*PrePrepareMsg, error)
	PrePrepare(prePrepareMsg *PrePrepareMsg) (*VoteMsg, error)
	Prepare(prepareMsg *VoteMsg) (*VoteMsg, error)
	Commit(commitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error)
}
