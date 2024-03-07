package main

import "log"

func main() {
	svc := InitService(123)
	err := svc.Register()
	if err != nil {
		log.Fatalf("service error: failed to init register %v", err)
	}
}
