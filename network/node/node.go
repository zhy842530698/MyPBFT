package node

import "MyPBFT_1/consensus"

type network_pbft interface {
	broadcast(msg interface{},path string)map[string]error
	reply(msg *consensus.ReplyMsg) error
	getReq(reqMsg *consensus.RequestMsg) error
	getPrePrepare(prePrepareMsg *consensus.PrePrepareMsg)error
	getPrepare(prepareMsg *consensus.VoteMsg) error
	getCommit(commitMsg *consensus.VoteMsg) error
	getReply(replyMsg *consensus.ReplyMsg)
	resolveRequestMsg(msgs []*consensus.RequestMsg) []error
	resolvePrePrepareMsg(msgs []*consensus.PrePrepareMsg)[]error
	resolvePrepareMsg(msgs []*consensus.VoteMsg) []error
	resolveCommitMsg(msgs []*consensus.VoteMsg) []error
	dispatchMsg()
	routeMsg(msg interface{}) []error
	routeMsgWhenAlarmed() []error
	alarmToDispatcher()
	resolveMsg()
}

