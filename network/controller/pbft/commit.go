package controller

import (
	//"MyPBFT_1/network/handler"
	handler "MyPBFT_1/network/handler/pbft"
	idl2 "MyPBFT_1/network/idl/pbft"
)



type PBFTCommitController struct {
}
func (c * PBFTCommitController) GenIdl() interface{} {
	return idl2.NewPBFTCommitIDL()
}

// Do ...
func (c *PBFTCommitController) Do(req interface{}) interface{} {
	var commit idl2.PBFTCommit
	r := req.(*idl2.PBFTCommit)
	commit.NodeID=r.NodeID
	commit.SequenceID=r.SequenceID
	commit.Digest=r.Digest
	commit.ViewID=r.ViewID
	commit.MsgType=r.MsgType
	commit.Sign=r.Sign
	commit.PVHash = r.PVHash
	//return handler.JoinToPeers(r)
	return handler.DeliverCommit(commit)
}