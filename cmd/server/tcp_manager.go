package main

import (
	"VulpesNet/pkg"
	"VulpesNet/utils"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TCPManager struct {
	ControlConn sync.Map
	WaitingList sync.Map
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

func (m *TCPManager) HandleRequest(c utils.Codec, header *utils.Header) {
	switch header.MagicNumber {
	case utils.NewSSH:
		var req pkg.SSHRequest
		err := c.ReadBody(&req)
		if err != nil {
			log.Println("server error: invalid request body ", req)
			return
		}
		codec, ok := m.ControlConn.Load(req.Data)
		if !ok {
			m.SendResponse(c, header, pkg.SSHResponse{Code: 1, Msg: "server error: not found requested service"})
			return
		}
		ch := make(chan utils.Codec)
		m.WaitingList.Store(req.ServiceId, ch)

		m.BuildReverseConn(codec.(utils.Codec), req.ServiceId)
		select {
		case sc := <-ch:
			go m.ProxySSH(c, sc)
			m.SendResponse(c, header, pkg.SSHResponse{Code: 0})
		case <-time.Tick(10 * time.Second):
			log.Println("server error: service reverse connection time out")
			return
		}
	case utils.RegisterService:
		var req pkg.SSHRequest
		err := c.ReadBody(&req)
		if err != nil {
			log.Println("server error: invalid request body ", req)
			return
		}
		m.ControlConn.Store(req.ServiceId, c)
	case utils.ReverseConn:
		var req pkg.SSHRequest
		err := c.ReadBody(&req)
		if err != nil {
			log.Println("server error: invalid request body ", req)
			return
		}
		some, ok := m.WaitingList.Load(req.Data)
		if !ok {
			log.Println("server error: failed to build reverse connection")
			return
		}
		ch := some.(chan utils.Codec)
		ch <- c
	}
}

func (m *TCPManager) ServeReverseCodec(codex utils.Codec) {
	time.Sleep(15 * time.Minute)
}

func (m *TCPManager) BuildReverseConn(c utils.Codec, serviceId uint64) {
	h := utils.Header{MagicNumber: utils.ReverseConn}
	b := pkg.SSHResponse{
		Code: 0,
		Data: serviceId,
	}
	_ = c.Write(&h, b)
}

func (m *TCPManager) ProxySSH(client utils.Codec, service utils.Codec) {
	log.Println("proxy...")
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
