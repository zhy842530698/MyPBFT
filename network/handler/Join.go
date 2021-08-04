package handler

import (
	"MyPBFT_1/network/common"
	"MyPBFT_1/network/idl"
	"fmt"
)

func JoinToPeers(request *idl.JoinRequest)*idl.JoinResponse {

	resp:=idl.NewJoinResponseIDL()
	common.DataQueue<-request
	err:= <- common.RespQueue
	fmt.Println("err:"+err.Error())
	resp.Errno = err.No
	resp.Msg = err.Msg
	return resp

}

//func
