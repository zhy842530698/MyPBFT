package controller

import (
	//"MyPBFT_1/network/handler"
	handler "MyPBFT_1/network/handler/pbft"
	idl2 "MyPBFT_1/network/idl/pbft"
)



type PBFTPrepareController struct {
}
func (c * PBFTPrepareController) GenIdl() interface{} {
	return idl2.NewPBFTPrepareIDL()
}

// Do ...
func (c *PBFTPrepareController) Do(req interface{}) interface{} {
	var prepare idl2.PBFTPrepare
	r := req.(*idl2.PBFTPrepare)
	prepare.ViewID=r.ViewID
	prepare.Digest=r.Digest
	prepare.SequenceID=r.SequenceID
	prepare.NodeID=r.NodeID
	prepare.MsgType=r.MsgType
	prepare.Sign = r.Sign
	prepare.PVHash=r.PVHash
	//return handler.JoinToPeers(r)
	return handler.DeliverPrepare(prepare)
}