package node

import (
	"MyPBFT_1/consensus"
	"MyPBFT_1/network/common"
	"MyPBFT_1/network/idl"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)
const ResolvingTimeDuration = time.Millisecond *1000 // 1 second.
const ViewID = "10024"
const Primay = "192.168.0.198:10024"
const ClientURL = "192.168.0.163:8080"

func NewNode(nodeID ,URL ,public string) *Node  {
	peers:=make(map[string]PubNode,4)
	peers["10024"]=PubNode{URL:"192.168.0.165:10024",Publickey: "./RSA/10024/public.pem"}
	peers["10025"]=PubNode{URL:"192.168.0.165:10025",Publickey: "./RSA/10025/public.pem"}
	peers["10026"]=PubNode{URL:"192.168.0.165:10026",Publickey: "./RSA/10026/public.pem"}
	peers["10027"]=PubNode{URL:"192.168.0.165:10027",Publickey: "./RSA/10027/public.pem"}
	node :=&Node{
		NodeID: nodeID,
		URL: URL,
		Peers: peers,
		View:&View{
			ID: ViewID,
			URL: Primay,
		},
		CurrentState: nil,
		CommittedMsgs: make([]*consensus.RequestMsg,0),
		MsgBuffer: &MsgBuffer{
			ReqMsgs: make([]*consensus.RequestMsg,0),
			PrePrepareMsgs: make([]*consensus.PrePrepareMsg,0),
			PrepareMsgs: make([]*consensus.VoteMsg,0),
		},
		MsgEntrance: make(chan interface{}),
		MsgDelivery: make(chan interface{}),
		Alarm: make(chan bool),
		Key: &Key{Public: public,private: strings.Replace(public,"public","private",-1)},
	}

	//start message dispatcher
	go node.dispatchMsg()
//	start alarm trigger
	go node.alarmToDispatcher()
//  resolve Msg
	go node.resolveMsg()
	return node

}

func (node *Node) broadcast(msg interface{}, path string) map[string]error {
	errorMap:=make(map[string]error)
	for ID,pubnode:=range node.Peers{
		if pubnode.URL==node.URL {
			continue
		}
		jsonMsg,err:=json.Marshal(msg)
		if err!=nil{
			errorMap[ID] = err
			continue
		}
		fmt.Println("******"+node.URL+" to "+pubnode.URL+path+"********")
		Send(pubnode.URL+path,jsonMsg)
		//time.Sleep(1000)
	}
	return errorMap
}



func (node *Node) getReq(reqMsg *consensus.RequestMsg) error {
	err:=node.createStateForNewConsensus()
	if err!=nil {
		return err
	}
	prepreparemsg,err:=node.CurrentState.StartConsensus(reqMsg)
	fmt.Println(prepreparemsg)
	//strconv.FormatInt 转String

	SignStr:=prepreparemsg.ViewID+strconv.FormatInt(prepreparemsg.SequenceID,10)

	prepreparemsg.Sign=common.GetSign([]byte(SignStr),node.Key.private)

	fmt.Printf("length is %d \n ",len(prepreparemsg.Sign))
	if err!=nil{
		return err
	}
	if prepreparemsg!=nil{
		node.broadcast(prepreparemsg,"/preprepare")
		fmt.Println("PrepreMsg is Done")
	}
	return nil
}

func (node *Node) getPrePrepare(prePrepareMsg *consensus.PrePrepareMsg) error {

	verfitymsg:=prePrepareMsg.ViewID+strconv.FormatInt(prePrepareMsg.SequenceID,10)
	if !common.VerifySign([]byte(verfitymsg),prePrepareMsg.Sign,node.Peers[prePrepareMsg.ViewID].Publickey){
			return errors.New("消息伪造")
	}

	err:=node.createStateForNewConsensus()
	if err!=nil{
		return err
	}
	pareMsg,err:=node.CurrentState.PrePrepare(prePrepareMsg)

	if err !=nil{
		return err
	}
	if pareMsg!=nil{
		Signstr:=node.NodeID+strconv.FormatInt(pareMsg.SequenceID,10)
		pareMsg.Sign=common.GetSign([]byte(Signstr),node.Key.private)
		pareMsg.NodeID = node.NodeID
		node.broadcast(pareMsg,"/prepare")
	}
	return nil
}

func (node *Node) getPrepare(prepareMsg *consensus.VoteMsg) error {
	fmt.Printf("Hash is %s \n",prepareMsg.Digest)
	verfitymsg:=prepareMsg.NodeID+strconv.FormatInt(prepareMsg.SequenceID,10)
	if !common.VerifySign([]byte(verfitymsg),prepareMsg.Sign,node.Peers[prepareMsg.NodeID].Publickey){
		return errors.New("消息伪造")
	}
	commitMsg,err:=node.CurrentState.Prepare(prepareMsg)
	if err!=nil{
		return err
	}
	if commitMsg != nil{
		Signstr:=node.NodeID+strconv.FormatInt(commitMsg.SequenceID,10)
		commitMsg.Sign=common.GetSign([]byte(Signstr),node.Key.private)
		commitMsg.NodeID = node.NodeID
		length:= len(node.Blocks)
		if length==0 {
			commitMsg.PVHash="0"
		}else {
			commitMsg.PVHash=node.Blocks[length-1].Hash
		}
		node.broadcast(commitMsg,"/commit")
	}
	return nil
}

func (node *Node) getCommit(commitMsg *consensus.VoteMsg) error {
	//验证消息
	fmt.Printf("Hash is %s \n",commitMsg.Digest)
	verfitymsg:=commitMsg.NodeID+strconv.FormatInt(commitMsg.SequenceID,10)
	if !common.VerifySign([]byte(verfitymsg),commitMsg.Sign,node.Peers[commitMsg.NodeID].Publickey){
		return errors.New("消息伪造")
	}
	replyMsg,commsg,error:=node.CurrentState.Commit(commitMsg)
	if error!=nil{
		return error
	}
	if replyMsg!=nil && commsg!=nil {
		replyMsg.NodeID  =node.NodeID
		Signstr:=node.NodeID+strconv.FormatInt(replyMsg.Timestamp,10)
		replyMsg.Sign=common.GetSign([]byte(Signstr),node.Key.private)

		block := node.GenerateBlocks(node.CurrentState.MsgLogs.CommitMsgs, commsg,commitMsg)
		if len(node.Blocks)==0 {
			block.PvHash="0"
			node.CommittedMsgs = append(node.CommittedMsgs,commsg)
			node.Blocks=append(node.Blocks,block)
		}else if len(node.Blocks)>0 &&node.Checkduplicate(block,node.Blocks[len(node.Blocks)-1].Hash) {
			fmt.Println("添加尾部")
			node.CommittedMsgs = append(node.CommittedMsgs,commsg)
			node.Blocks=append(node.Blocks,block)
		}
		node.reply(replyMsg)

	}
	return nil
}
func (node *Node) reply(msg *consensus.ReplyMsg) error {
	////持久化
	fmt.Println("blocks")
	for _,block:=range node.Blocks{
		fmt.Println(block)
	}
	jsonMsg,_:=json.Marshal(msg)
	fmt.Println("发送给客户端")
	Send(ClientURL+"/reply",jsonMsg)


	return nil
}
func (node *Node) GetReply(replyMsg *consensus.ReplyMsg) {
	fmt.Print("currentstate: ")
	fmt.Println(node.CurrentState.MsgLogs.ReqMsg)
	fmt.Println(node.CurrentState.MsgLogs.CommitMsgs)
	fmt.Printf("Result:%s by %s\n",replyMsg.Result,replyMsg.NodeID)
}


func (node *Node) resolveRequestMsg(msgs []*consensus.RequestMsg) []error {
	errs:=make([]error,0)
	//resolve Message
	for _,reqMsg:=range msgs{
		err:=node.getReq(reqMsg)
		if err!=nil{
			errs=append(errs,err)
		}
	}
	if len(errs)!=0{
		return errs
	}
	return  nil
}

func (node *Node) resolvePrePrepareMsg(msgs []*consensus.PrePrepareMsg) []error {
	errs:=make([]error,0)
	//resolve Message
	for _,premsg:=range msgs{
		err:=node.getPrePrepare(premsg)
		if err!=nil{
			errs=append(errs,err)
		}
	}
	if len(errs)!=0{
		return errs
	}
	return  nil
}

func (node *Node) resolvePrepareMsg(msgs []*consensus.VoteMsg) []error {
	errs:=make([]error,0)
	//resolve Message
	for _,premsg:=range msgs{
		err:=node.getPrepare(premsg)
		if err!=nil{
			errs=append(errs,err)
		}
	}
	if len(errs)!=0{
		return errs
	}
	return  nil
}

func (node *Node) resolveCommitMsg(msgs []*consensus.VoteMsg) []error {
	errs := make([]error, 0)

	// Resolve messages
	for _, commitMsg := range msgs {
		err := node.getCommit(commitMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}
//dispatchMsg()
//routeMsg(msg interface{}) []error
//routeMsgWhenAlarmed() []error
//alarmToDispatcher()

func (node *Node)dispatchMsg()  {
	for  {
		select {
		 case msg:=<-node.MsgEntrance:
		 	//fmt.Println("进入分发")
		 	err:=node.routeMsg(msg)
			 if err!=nil {
				 fmt.Println(err)
			 }
		case msg:=<-common.DataQueue:
			//fmt.Println("进入判决")
			err:=node.routeMsgX(msg)
			if err!=nil {
				common.RespQueue<-err
			}
		case <-node.Alarm:
			//fmt.Println("进入Alarm")
			err := node.routeMsgWhenAlarmed()
			if err !=nil{
				fmt.Println(err)
			}
		}
	}
}
func (node *Node)routeMsgX(msg interface{})*common.Myerror{
	//fmt.Println(reflect.TypeOf(msg))
	switch msg.(type) {
	case *idl.JoinRequest:

		err:=node.Join(msg.(*idl.JoinRequest))
		if err!=nil{
			return err
		}
	}
	return nil
}

func(node *Node)routeMsg(msg interface{}) []error  {
	//fmt.Println(msg.(reflect.Type))
	switch msg.(type) {
	case *consensus.RequestMsg:
		node.CurrentState = nil
		if node.CurrentState==nil{

			msgs:=make([]*consensus.RequestMsg,len(node.MsgBuffer.ReqMsgs))
			copy(msgs,node.MsgBuffer.ReqMsgs)
			msgs=append(msgs,msg.(*consensus.RequestMsg))
			node.MsgBuffer.ReqMsgs = make([]*consensus.RequestMsg,0)
			node.MsgDelivery <- msgs
		}else {
			node.MsgBuffer.ReqMsgs=append(node.MsgBuffer.ReqMsgs,msg.(*consensus.RequestMsg))
			//fmt.Println(node.MsgBuffer.ReqMsgs)
		}

	case *consensus.PrePrepareMsg:
		node.CurrentState = nil
		if node.CurrentState==nil{
			msgs :=make([]*consensus.PrePrepareMsg,len(node.MsgBuffer.PrePrepareMsgs))
			copy(msgs,node.MsgBuffer.PrePrepareMsgs)
			msgs=append(msgs,msg.(*consensus.PrePrepareMsg))
			node.MsgBuffer.PrePrepareMsgs=make([]*consensus.PrePrepareMsg,0)
			node.MsgDelivery <- msgs
		}else{
			node.MsgBuffer.PrePrepareMsgs = append(node.MsgBuffer.PrePrepareMsgs,msg.(*consensus.PrePrepareMsg))
		}
	case *consensus.VoteMsg:
		//fmt.Println(msg)
		if msg.(*consensus.VoteMsg).MsgType==consensus.PrepareMsg{
			if node.CurrentState == nil || node.CurrentState.CurrentStage != consensus.PrePrepared {
				node.MsgBuffer.PrepareMsgs = append(node.MsgBuffer.PrepareMsgs, msg.(*consensus.VoteMsg))
			} else {
				// Copy buffered messages first.
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.PrepareMsgs))
				copy(msgs, node.MsgBuffer.PrepareMsgs)

				// Append a newly arrived message.
				msgs = append(msgs, msg.(*consensus.VoteMsg))

				// Empty the buffer.
				node.MsgBuffer.PrepareMsgs = make([]*consensus.VoteMsg, 0)

				// Send messages.
				node.MsgDelivery <- msgs
			}
		}else if msg.(*consensus.VoteMsg).MsgType == consensus.CommitMsg {
			if node.CurrentState == nil || node.CurrentState.CurrentStage != consensus.Prepared {
				fmt.Println("分发commitMsg")
				node.MsgBuffer.CommitMsgs = append(node.MsgBuffer.CommitMsgs, msg.(*consensus.VoteMsg))
			} else {
				// Copy buffered messages first.
				fmt.Println("添加其他人的CommitMsg")
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.CommitMsgs))
				copy(msgs, node.MsgBuffer.CommitMsgs)

				// Append a newly arrived message.
				msgs = append(msgs, msg.(*consensus.VoteMsg))

				// Empty the buffer.
				node.MsgBuffer.CommitMsgs = make([]*consensus.VoteMsg, 0)

				// Send messages.

				node.MsgDelivery <- msgs
			}
		}
	case *consensus.ReplyMsg:


	}
	return nil
}
func (node *Node)routeMsgWhenAlarmed() []error  {
	if node.CurrentState == nil {
		// Check ReqMsgs, send them.
		if len(node.MsgBuffer.ReqMsgs) != 0 {
			msgs := make([]*consensus.RequestMsg, len(node.MsgBuffer.ReqMsgs))
			copy(msgs, node.MsgBuffer.ReqMsgs)

			node.MsgDelivery <- msgs
		}

		// Check PrePrepareMsgs, send them.
		if len(node.MsgBuffer.PrePrepareMsgs) != 0 {
			msgs := make([]*consensus.PrePrepareMsg, len(node.MsgBuffer.PrePrepareMsgs))
			copy(msgs, node.MsgBuffer.PrePrepareMsgs)

			node.MsgDelivery <- msgs
		}
	} else {
		switch node.CurrentState.CurrentStage {
		case consensus.PrePrepared:
			// Check PrepareMsgs, send them.
			if len(node.MsgBuffer.PrepareMsgs) != 0 {
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.PrepareMsgs))
				copy(msgs, node.MsgBuffer.PrepareMsgs)

				node.MsgDelivery <- msgs
			}
		case consensus.Prepared:
			// Check CommitMsgs, send them.
			if len(node.MsgBuffer.CommitMsgs) != 0 {
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.CommitMsgs))
				copy(msgs, node.MsgBuffer.CommitMsgs)
				node.MsgDelivery <- msgs
			}
		}
	}

	return nil
}
func (node *Node)alarmToDispatcher()  {
	for {
		time.Sleep(ResolvingTimeDuration)
		node.Alarm <- true
	}
}

func (node *Node) resolveMsg() {
	for{
		msgs:=<-node.MsgDelivery
		switch msgs.(type) {
		case []*consensus.RequestMsg:
			//fmt.Println("解析RequestMsg")
			errs:=node.resolveRequestMsg(msgs.([]*consensus.RequestMsg))
			if len(errs)!=0 {
				for _,err := range errs{
					fmt.Println(err)
				}
			}
		case []*consensus.PrePrepareMsg:
			errs := node.resolvePrePrepareMsg(msgs.([]*consensus.PrePrepareMsg))
			if len(errs) != 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				// TODO: send err to ErrorChannel
			}
		case []*consensus.VoteMsg:
			voteMsgs := msgs.([]*consensus.VoteMsg)
			if len(voteMsgs) == 0 {
				break
			}
			if voteMsgs[0].MsgType == consensus.PrepareMsg {
				errs := node.resolvePrepareMsg(voteMsgs)
				if len(errs) != 0 {
					for _, err := range errs {
						fmt.Println(err)
					}
					// TODO: send err to ErrorChannel
				}
			} else if voteMsgs[0].MsgType == consensus.CommitMsg {
				errs := node.resolveCommitMsg(voteMsgs)
				if len(errs) != 0 {
					for _, err := range errs {
						fmt.Println(err)
					}
					// TODO: send err to ErrorChannel
				}
			}


		}
	}

}
func(node *Node) createStateForNewConsensus() error{
	if node.CurrentState!=nil{
		return errors.New("another consensus is going")
	}
	var lastSequenceID int64
	if len(node.CommittedMsgs)==0{
		lastSequenceID=-1;
	}else{
		lastSequenceID = node.CommittedMsgs[len(node.CommittedMsgs)-1].SequenceID
	}
	node.CurrentState = consensus.CreateState(node.View.ID,lastSequenceID)
	return nil
}
func(node *Node)Checkduplicate(block *Block,Hash string)bool{
	if block.PvHash!=Hash{
		fmt.Println("pvhash is :")
		fmt.Println(block.PvHash)
		fmt.Println("---------------")
		fmt.Println("hash is :")
		fmt.Println(block.Hash)
		return false
	}
	return true
}






