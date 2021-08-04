package consensus


type RequestMsg struct {
	Timestamp  int64  `json:"timestamp"`
	ClientID   string `json:"clientID"`
	Operation  string `json:"operation"`
	SequenceID int64  `json:"sequenceID"`
}

type ReplyMsg struct {
	ViewID    string `json:"viewID"`
	Timestamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	NodeID    string `json:"nodeID"`
	Result    string `json:"result"`
	Sign      []byte  `json:"sign"`
	Ip        string  	`json:ip`
}

type PrePrepareMsg struct {

	ViewID     string       `json:"viewID"`
	SequenceID int64       `json:"sequenceID"`
	Digest     string      `json:"digest"`
	RequestMsg *RequestMsg `json:"requestMsg"`
	Sign        []byte   	`json:"sign"`
}

type VoteMsg struct {
	ViewID     string  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     string `json:"digest"`
	NodeID     string `json:"nodeID"`
	Sign        []byte  `json:"sign"`
	PVHash     string  `json:"pvhash"`
	MsgType           `json:"msgType"`
}

type MsgType int
const (
	PrepareMsg MsgType = iota
	CommitMsg
)

