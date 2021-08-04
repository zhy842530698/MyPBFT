package handler

import (
	"MyPBFT_1/consensus"
	"MyPBFT_1/network/common"
	idl2 "MyPBFT_1/network/idl/pbft"
	"github.com/mapstructure"
)

func DeliverRequest(request idl2.PBFTRequest)interface{}  {
	var msg consensus.RequestMsg
	//var Err *common.Myerror
	resp:=idl2.NewPBFTRequestResponseIDL()
	//fmt.Println("----------")
	//fmt.Println(request)
	mapre:=common.StructAtoB(request,msg)
	//fmt.Println(mapre)
	if err:=mapstructure.Decode(mapre,&msg);err!=nil{
		//fmt.Println(err)
		resp.Errno=-1
		resp.Msg="数据解析错误"
		return resp
	}
	resp.Errno=0
	resp.Msg="request阶段"
	common.HttpQueue<-&msg
	return resp
}
