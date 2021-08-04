package controller

import (
	//"MyPBFT_1/network/handler"
	handler "MyPBFT_1/network/handler/pbft"
	idl2 "MyPBFT_1/network/idl/pbft"
)



type PBFTPrePrepareController struct {
}
func (c *PBFTPrePrepareController) GenIdl() interface{} {
	return idl2.NewPBFTPrePrepareIDL()
}

// Do ...
func (c *PBFTPrePrepareController) Do(req interface{}) interface{} {
	var prepre idl2.PBFTPreprepare
	r := req.(*idl2.PBFTPreprepare)
	prepre.ViewID=r.ViewID
	prepre.SequenceID=r.SequenceID
	prepre.RequestMsg=r.RequestMsg
	prepre.Digest=r.Digest
	prepre.Sign=r.Sign

	//return handler.JoinToPeers(r)
	return handler.DeliverPreprepare(prepre)
}