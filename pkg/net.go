package pkg

type SSHRequest struct {
	ServiceId uint64
	Token     string
	Data      uint64
}

type SSHResponse struct {
	Code int
	Msg  string
	Data uint64
}
