package handler

import (
	"MyPBFT_1/consensus"
	"MyPBFT_1/network/common"
	idl2 "MyPBFT_1/network/idl/pbft"
	"github.com/mapstructure"
)

func DeliverCommit(request idl2.PBFTCommit)interface{}  {
	var msg consensus.VoteMsg
	//var Err *common.Myerror
	resp:=idl2.NewPBFTCommitResponseIDL()
	//fmt.Println(request)
	mapre:=common.StructAtoB(request,msg)
	if err:=mapstructure.Decode(mapre,&msg);err!=nil{
		//fmt.Println(err)
		resp.Errno=-1
		resp.Msg="数据解析错误"
		return resp
	}
	resp.Errno=0
	resp.Msg="Commit阶段成功"
	common.HttpQueue<-&msg
	return resp
}
