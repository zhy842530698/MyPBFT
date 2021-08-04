package idl

import "MyPBFT_1/consensus"

// CRequest request struct
type PBFTPrepare struct {
	ViewID     string  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     string `json:"digest"`
	NodeID     string `json:"nodeID"`
	Sign        []byte  	`json:"sign"`
	PVHash     string  `json:"pvhash"`
	consensus.MsgType           `json:"msgType"`
}

// NewCRequestIDL ...
func NewPBFTPrepareIDL() *PBFTPrepare {
	return &PBFTPrepare{}
}

// CResponse response struct
type PBFTPrepareResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewCResponseIDL ...
func NewPBFTPrepareResponseIDL() *PBFTPrepareResponse {
	return &PBFTPrepareResponse{}
}

