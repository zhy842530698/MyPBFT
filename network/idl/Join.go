package idl
// CRequest request struct
type JoinRequest struct {
	Userid string `json:"userid"`
	Url string `json:"url"`
}

// NewCRequestIDL ...
func NewJoinRequestIDL() *JoinRequest {
	return &JoinRequest{}
}

// CResponse response struct
type JoinResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

// NewCResponseIDL ...
func NewJoinResponseIDL() *JoinResponse {
	return &JoinResponse{}
}

