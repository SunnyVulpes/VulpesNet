package service

import (
	"VulpesNet/server"
	"VulpesNet/utils"
	"log"
	"net"
)

type Service struct {
	serviceId uint64
	codec     utils.Codec
}

func InitService(serviceId uint64) *Service {
	conn, err := net.Dial("tcp", ":9000")
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
	b := server.SSHRequest{
		ServiceId: svc.serviceId,
		Token:     "tat",
	}
	return svc.codec.Write(&h, b)
}
