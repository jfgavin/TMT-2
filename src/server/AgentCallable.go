package server

import "github.com/jfgavin/TMT-2/src/agent"

// Functions explicitly callable by the agent via the ServerAPI will be placed here
func (serv *GameServer) GetEliminationCount() int {
	return serv.elims
}

func (serv *GameServer) RequestSacrifice(sacAg *agent.TMTAgent) {
	serv.sacrificeRequests = append(serv.sacrificeRequests, sacAg)
}
