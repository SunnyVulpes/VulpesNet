package client

import (
	"VulpesNet/server"
	"VulpesNet/utils"
	"io"
	"log"
	"net"
	"sync"
)

type Client struct {
	token   string
	sending sync.Mutex
	codec   utils.Codec
}

func InitClient() *Client {
	return &Client{
		token: "qaq",
	}
}

func (c *Client) DialSSH() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		log.Println("client error: dial to server failed ", err)
		return
	}
	c.codec = utils.InitJSONCodec(conn)

	c.sending.Lock()
	defer c.sending.Unlock()
	err = c.codec.Write(&utils.Header{
		MagicNumber: utils.NewSSH,
	}, &server.SSHRequest{
		ServiceId: 123,
		Token:     "qaq",
	})
	if err != nil {
		log.Println("client error: write to server failed", err)
	}

	go c.Receive()
}

func (c *Client) Receive() {
	var err error
	for {
		var header utils.Header
		err = c.codec.ReadHeader(&header)
		if err != nil {
			break
		}

		switch header.MagicNumber {
		case utils.NewSSH:
			response := server.SSHResponse{}
			err = c.codec.ReadBody(&response)
			if response.Code == 0 {
				go c.ProxySSH()
				return
			} else {
				log.Println("client: ssh request was rejected ", response.Msg)
			}
		}
	}
	log.Println("client error: receive failed ", err)
}

func (c *Client) ProxySSH() {
	l, err := net.Listen("tcp", ":7070")
	if err != nil {
		log.Fatalf("client fatal: listen local ssh failed")
		return
	}
	conn, _ := l.Accept()
	go io.Copy(conn, c.codec.GetConn())
	io.Copy(c.codec.GetConn(), conn)
}
