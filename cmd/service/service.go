package main

import (
	"VulpesNet/pkg"
	"VulpesNet/utils"
	"io"
	"log"
	"net"
	"sync"
)

type Service struct {
	serviceId   uint64
	codec       utils.Codec
	reverseConn sync.Map
}

func InitService(serviceId uint64) *Service {
	conn, err := net.Dial("tcp", utils.Host+":9000")
	if err != nil {
		log.Fatalf("service error: failed to init dial %v", err)
		return nil
	}

	codec := utils.InitJSONCodec(conn)
	svc := &Service{
		serviceId: serviceId,
		codec:     codec,
	}

	return svc
}

func (svc *Service) Register() error {
	h := utils.Header{
		MagicNumber: utils.RegisterService,
	}
	b := pkg.SSHRequest{
		ServiceId: svc.serviceId,
		Token:     "tat",
	}
	go svc.ReceiveCmd()
	err := svc.codec.Write(&h, b)
	if err != nil {
		return err
	}
	for {
	}
}

func (svc *Service) ReceiveCmd() {
	for {
		var header utils.Header
		err := svc.codec.ReadHeader(&header)
		if err != nil {
			log.Println("service error: read cmd header failed ", err)
			return
		}
		switch header.MagicNumber {
		case utils.ReverseConn:
			var response pkg.SSHResponse
			err = svc.codec.ReadBody(&response)
			if err != nil || response.Code != 0 {
				log.Println("service error: read cmd body failed ", err)
				return
			}
			go svc.DialReverseConn(response.Data)
		}
	}
}

func (svc *Service) DialReverseConn(serviceId uint64) {
	conn, err := net.Dial("tcp", utils.Host+":9000")
	if err != nil {
		log.Println("service error: dial reverse connection failed ", err)
		return
	}
	codec := utils.InitJSONCodec(conn)
	codec.Write(&utils.Header{
		MagicNumber: utils.ReverseConn,
	}, pkg.SSHRequest{ServiceId: svc.serviceId, Data: serviceId})

	svc.ProxySSH(codec)
}

func (svc *Service) ProxySSH(c utils.Codec) {
	conn, err := net.Dial("tcp", ":22")
	if err != nil {
		log.Println("service error: dial ssh failed ", err)
	}
	go io.Copy(conn, c.GetConn())
	io.Copy(c.GetConn(), conn)
	log.Println("proxy...")
}
