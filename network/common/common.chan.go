package common

var (
	DataQueue  chan interface{}
	RespQueue  chan  *Myerror
	//RequestQueue chan *consensus.RequestMsg
	//PreprepareQueue chan *consensus.PrePrepareMsg
	//PrepareQueue chan *consensus.VoteMsg
	//CommitQueue  chan *consensus.VoteMsg
	//ReplyQueue   chan *consensus.ReplyMsg
	HttpQueue chan interface{}
)