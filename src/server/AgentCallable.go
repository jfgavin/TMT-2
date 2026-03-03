package server

import "github.com/jfgavin/TMT-2/src/agent"

// Functions explicitly callable by the agent via the ServerAPI will be placed here

func (serv *GameServer) GetEliminationCount() int {
	return serv.elims
}

// Agents can pass a pointer to themselves when requesting their sacrifice
// It is important to consider security of this method
// Whilst perhaps impossible, an agent which can reference other agents could maliciously request their deaths
// By using a pointer, agents should not have access to each other's pointers, so should be secure, unlike e.g. Name
func (serv *GameServer) RequestSacrifice(sacAg *agent.TMTAgent) {
	serv.SacrificeAgent(sacAg)
}
