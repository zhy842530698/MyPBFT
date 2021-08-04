package node

import (
	"MyPBFT_1/consensus"
	"MyPBFT_1/network/common"
	"MyPBFT_1/network/idl"
	"crypto/sha256"
	"fmt"
	"github.com/cbergoon/merkletree"
	"strconv"
	"strings"
	"sync/atomic"
)

const split =":"
func(node *Node)Join(request *idl.JoinRequest)*common.Myerror{
	var err *common.Myerror
	url := request.Url
	userid := request.Userid
	arr:=strings.Split(url,split)
	if len(arr)<=1 {
		err=&common.Myerror{Msg: "请指定端口号和正确的url",No: -3}
		return err
	}
	if node.Peers[userid].URL=="" {
		str:=common.ValidIPAddress(arr[0])
		_,error:=strconv.Atoi(arr[1])
		if str=="IPv4"&& error==nil{
			//node.Peers[userid]=PubNode{URL:url,Publickey: "" }
		}else{
			err=&common.Myerror{"无效的url地址",-1}
			return err
		}
	}else {
		err:=&common.Myerror{"节点id重复请重新更换",-2}
		return err
	}
	err=&common.Myerror{"成功",0}
	return err
}
func (t TreeNode) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.val)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
func (t TreeNode) Equals(other merkletree.Content) (bool, error) {
	return t.val == other.(TreeNode).val, nil
}
func(node *Node) GenerateBlocks(votesmsgs map[string]*consensus.VoteMsg,reqmsg *consensus.RequestMsg,commsg *consensus.VoteMsg)*Block{
	var list []merkletree.Content
	list =make([]merkletree.Content,0)
	for _,votemsg:=range votesmsgs {
		//fmt.Println(votemsg)
		list=append(list, TreeNode{votemsg.NodeID})
	}
	fmt.Println("commsg------>>>>>>>>:")
	fmt.Println(commsg)
	metaData:=string(reqmsg.Timestamp)+string(reqmsg.SequenceID)+commsg.PVHash+reqmsg.Operation
	hash, _:= HashwithDifficulty([]byte(metaData), 3)
	t,err:=merkletree.NewTree(list)
	if err!=nil  {
		fmt.Printf("merkletree Error is ",err)
	}
	return &Block{
		Req: reqmsg,
		Tree: t,
		PvHash:commsg.PVHash,
		Hash: fmt.Sprintf("%x", hash),
	}

}
func HashwithDifficulty(data []byte, d int) (result [32]byte, nonce int64) {
	var stop int32
	for nonce = 1; ; nonce++ {
		if atomic.LoadInt32(&stop) == 1 {
			return result, 0
		}
		str := strconv.FormatInt(nonce, 10)
		b := append(data, []byte(str)...)
		result = sha256.Sum256(b)
		if difficulty(result, d) {
			return result, nonce
		}
	}
	return
}
func difficulty(hash [32]byte, d int) bool {
	dn := d / 2
	sn := d % 2
	for i := 0; i < dn; i++ {
		if hash[i] != 0x00 {
			return false
		}
	}
	if sn != 0 {
		if hash[dn*2+1] > 0x0f {
			return false
		}
	}
	return true
}

