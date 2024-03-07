package main

import (
	main2 "VulpesNet/cmd/client"
	"VulpesNet/cmd/server"
	"VulpesNet/cmd/service"
	"log"
	"time"
)

func Server() {
	s := server.InitServer()
	s.Run()
}

func Service() {
	svc := service.InitService(123)
	err := svc.Register()
	if err != nil {
		log.Fatalf("service error: failed to init register %v", err)
	}
	time.Sleep(15 * time.Minute)
}

func Client() {
	time.Sleep(2 * time.Second)
	c := main2.InitClient()
	c.DialSSH()
	//time.Sleep(2 * time.Second)
	//conn, err := net.Dial("tcp", ":9090")
	//if err != nil {
	//	log.Println(err)
	//}
	//_, err = conn.Write([]byte("test"))
	//if err != nil {
	//	log.Println(err)
	//}

	time.Sleep(15 * time.Minute)
}

func main() {
	go Server()
	go Service()
	Client()
}
