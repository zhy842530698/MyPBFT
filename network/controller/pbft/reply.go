package controller

import (
	//"MyPBFT_1/network/handler"
	handler "MyPBFT_1/network/handler/pbft"
	idl2 "MyPBFT_1/network/idl/pbft"
)



type PBFTReplyController struct {
}
func (c *PBFTReplyController) GenIdl() interface{} {
	return idl2.NewPBFTReplyIDL()
}

// Do ...
func (c *PBFTReplyController) Do(req interface{}) interface{} {
	var reply idl2.PBFTReply
	r := req.(*idl2.PBFTReply)
	reply.ClientID=r.ClientID
	reply.Timestamp=r.Timestamp
	reply.NodeID=r.NodeID
	reply.Result=r.Result
	reply.ViewID=r.ViewID
	reply.Sign = r.Sign
	//return handler.JoinToPeers(r)
	return handler.DeliverReply(reply)
}