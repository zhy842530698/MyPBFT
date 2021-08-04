package idl
// CRequest request struct
type PBFTReply struct {
	ViewID    string `json:"viewID"`
	Timestamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	NodeID    string `json:"nodeID"`
	Result    string `json:"result"`
	Sign        []byte  	`json:"sign"`
}

// NewCRequestIDL ...
func NewPBFTReplyIDL() *PBFTReply {
	return &PBFTReply{}
}

// CResponse response struct
type PBFTReplyResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewCResponseIDL ...
func NewPBFTReplyResponseIDL() *PBFTReplyResponse {
	return &PBFTReplyResponse{}
}

