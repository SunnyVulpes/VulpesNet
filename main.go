package main

import (
	"VulpesNet/client"
	"VulpesNet/server"
	"VulpesNet/service"
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
	time.Sleep(5 * time.Second)
	c := client.InitClient()
	c.DialSSH()
	time.Sleep(15 * time.Minute)
}

func main() {
	go Server()
	go Service()
	Client()
}
