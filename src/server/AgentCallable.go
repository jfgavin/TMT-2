package server

// Functions explicitly callable by the agent via the ServerAPI will be placed here

func (serv *GameServer) GetEliminationCount() int {
	return serv.elims
}
