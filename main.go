package main

import (
	"github.com/jfgavin/TMT-2/src/config"
	gameServer "github.com/jfgavin/TMT-2/src/server"
)

// "go run ."
func main() {
	cfg := config.NewConfig()
	// create server from constructor
	serv := gameServer.NewGameServer(cfg)
	// begin simulator
	serv.Start()
}
