package main

import (
	"github.com/dapperlabs/flow-go-sdk/language/languageserver/server"
)

func main() {
	server.NewServer().Start()
}
