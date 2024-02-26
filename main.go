package main

import (
	"go-vault/api"
	"go-vault/internal/utility"
)

func main() {
	utility.InitLogger()
	api.NewAPIServer(":5464", "*").Run()
}
