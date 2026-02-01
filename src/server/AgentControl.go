package server

func (serv *GameServer) ElimDrainedAgents() {
	for _, ag := range serv.GetAgentMap() {
		if ag.GetEnergy() <= 0 {
			serv.RemoveAgent(ag)
		}
	}
}
