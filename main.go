package main

import gameServer "github.com/jfgavin/TMT-2/src/server"

// "go run ."
func main() {
	// create server from constructor
	serv := gameServer.MakeGameServer(2, 5)
	// toggle verbose logging of messaging stats
	serv.ReportMessagingDiagnostics()
	// begin simulator
	serv.Start()
}
