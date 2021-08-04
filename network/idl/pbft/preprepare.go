package idl

import "MyPBFT_1/consensus"

// CRequest request struct
type PBFTPreprepare struct {
	ViewID     string       `json:"viewID"`
	SequenceID int64       `json:"sequenceID"`
	Digest     string      `json:"digest"`
	RequestMsg *consensus.RequestMsg `json:"requestMsg"`
	Sign        []byte   	`json:"sign"`
}

// NewCRequestIDL ...
func NewPBFTPrePrepareIDL() *PBFTPreprepare {
	return &PBFTPreprepare{}
}

// CResponse response struct
type PBFTPreprepareResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewCResponseIDL ...
func NewPBFTPrePrepareResponseIDL() *PBFTPreprepareResponse {
	return &PBFTPreprepareResponse{}
}

