package main

import (
	"go-vault/api"
)

func main() {
	s := &api.APIServer{
		ListenerAddress: "4269",
		CORS:            "*",
	}
	s.Run()
}
