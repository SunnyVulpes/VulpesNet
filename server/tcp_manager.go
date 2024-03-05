package server

import (
	"VulpesNet/utils"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TCPManager struct {
	ReverseConn sync.Map
}

func (m *TCPManager) Run() {
	m.ListenClient()
}

func (m *TCPManager) ListenClient() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("server error: tcp manager listen client failed %v", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("server error: accept a new tcp connection failed %v", err)
		}
		go m.ServeConn(conn)
	}
}

func (m *TCPManager) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()
	c := utils.InitJSONCodec(conn)
	m.ServeCodec(c)
}

func (m *TCPManager) ServeCodec(c utils.Codec) {
	var header utils.Header
	err := c.ReadHeader(&header)
	if err != nil {
		log.Println("server err: read header failed", err)
		return
	}

	m.HandleRequest(c, &header)
}

type SSHRequest struct {
	ServiceId uint64
	Token     string
}

type SSHResponse struct {
	Code int
	Msg  string
}

func (m *TCPManager) HandleRequest(c utils.Codec, header *utils.Header) {
	switch header.MagicNumber {
	case utils.NewSSH:
		var req SSHRequest
		err := c.ReadBody(&req)
		if err != nil {
			log.Println("server error: invalid request body ", req)
			return
		}
		serviceCodec, ok := m.ReverseConn.Load(req.ServiceId)
		if !ok {
			m.SendResponse(c, header, SSHResponse{Code: 1, Msg: "server error: not found requested service"})
			return
		}
		m.SendResponse(c, header, SSHResponse{Code: 0})
		m.ProxySSH(c, serviceCodec.(utils.Codec))
	case utils.RegisterService:
		var req SSHRequest
		err := c.ReadBody(&req)
		if err != nil {
			log.Println("server error: invalid request body ", req)
			return
		}
		m.ReverseConn.Store(req.ServiceId, c)
		go m.ServeReverseCodec(c)
	}
}

func (m *TCPManager) ServeReverseCodec(codex utils.Codec) {
	time.Sleep(15 * time.Minute)
}

func (m *TCPManager) ProxySSH(client utils.Codec, service utils.Codec) {
	go func() {
		_, _ = io.Copy(client.GetConn(), service.GetConn())
	}()
	_, _ = io.Copy(service.GetConn(), client.GetConn())
}

func (m *TCPManager) SendResponse(c utils.Codec, header *utils.Header, body interface{}) {
	err := c.Write(header, body)
	if err != nil {
		log.Println("server error: send response failed ", err)
	}
}
