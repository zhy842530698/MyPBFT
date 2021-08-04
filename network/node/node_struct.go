package node

import (
	"MyPBFT_1/consensus"
	"github.com/cbergoon/merkletree"
)

type Node struct {
	NodeID        string
	URL           string
	Peers  		  map[string]PubNode
	View          *View
	CurrentState  *consensus.State
	CommittedMsgs []*consensus.RequestMsg // kinda block.
	MsgBuffer     *MsgBuffer
	MsgEntrance   chan interface{}
	MsgDelivery   chan interface{}
	Alarm         chan bool
	Key           *Key
	Blocks    	[]*Block

}
type Key struct {
	Public  string
	private string
}
type PubNode struct {
	URL       string
	Publickey string
}
type MsgBuffer struct {
	ReqMsgs        []*consensus.RequestMsg
	PrePrepareMsgs []*consensus.PrePrepareMsg
	PrepareMsgs    []*consensus.VoteMsg
	CommitMsgs     []*consensus.VoteMsg
}

type View struct {
	ID string
	URL string
}
type TreeNode struct {
	val string
}
type Block struct {
	 Req      *consensus.RequestMsg
	 Tree     *merkletree.MerkleTree
	 PvHash   string
	 Hash string
}

