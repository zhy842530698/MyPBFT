package idl
// CRequest request struct
type PBFTRequest struct {
	Timestamp  int64  `json:"timestamp"`
	ClientID   string `json:"clientID"`
	Operation  string `json:"operation"`
	SequenceID int64  `json:"sequenceID"`
}

// NewCRequestIDL ...
func NewPBFTRequestIDL() *PBFTRequest {
	return &PBFTRequest{}
}

// CResponse response struct
type PBFTRequestResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewCResponseIDL ...
func NewPBFTRequestResponseIDL() *PBFTRequestResponse {
	return &PBFTRequestResponse{}
}

