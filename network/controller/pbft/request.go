package controller

import (
	handler2 "MyPBFT_1/network/handler/pbft"
	idl2 "MyPBFT_1/network/idl/pbft"
)



type PBFTRequestController struct {
}
func (c *PBFTRequestController) GenIdl() interface{} {
	return idl2.NewPBFTRequestIDL()
}

// Do ...
func (c *PBFTRequestController) Do(req interface{}) interface{} {
	var request idl2.PBFTRequest
	r := req.(*idl2.PBFTRequest)
	request.SequenceID=r.SequenceID
	request.Timestamp=r.Timestamp
	request.ClientID=r.ClientID
	request.Operation= r.Operation
	return handler2.DeliverRequest(request)
}