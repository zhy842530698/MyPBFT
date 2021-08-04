package idl

import "MyPBFT_1/consensus"

// CRequest request struct
type PBFTCommit struct {
	ViewID     string  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     string `json:"digest"`
	NodeID     string `json:"nodeID"`
	Sign        []byte  	`json:"sign"`
	PVHash     string  `json:"pvhash"`
	consensus.MsgType           `json:"msgType"`
}

// NewCRequestIDL ...
func NewPBFTCommitIDL() *PBFTCommit {
	return &PBFTCommit{}
}

// CResponse response struct
type PBFTPBFTCommitResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewCResponseIDL ...
func NewPBFTCommitResponseIDL() *PBFTPBFTCommitResponse {
	return &PBFTPBFTCommitResponse{}
}


