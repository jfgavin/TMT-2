package server

func (serv *GameServer) DrainAgents() {
	for _, ag := range serv.GetAgentMap() {
		ag.ChangeEnergy(-1)
	}
}

func (serv *GameServer) ElimDrainedAgents() {
	for _, ag := range serv.GetAgentMap() {
		if ag.GetEnergy() < 0 {
			serv.RemoveAgent(ag)
		}
	}
}
