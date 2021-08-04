package controller

import (
	"MyPBFT_1/network/handler"
	"MyPBFT_1/network/idl"
)

type JoinController struct {
}
func (c *JoinController) GenIdl() interface{} {
return idl.NewJoinRequestIDL()
}

// Do ...
func (c *JoinController) Do(req interface{}) interface{} {
	r := req.(*idl.JoinRequest)
	return handler.JoinToPeers(r)
}