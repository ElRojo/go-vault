package main

import (
	"go-vault/api"
)

func main() {
	api.NewAPIServer(4269, "*").Run()
}
