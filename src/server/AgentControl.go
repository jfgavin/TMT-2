package server

import (
	"math/rand"

	"github.com/jfgavin/TMT-2/src/agent"
)

func (serv *GameServer) GetShuffledAgents() []agent.ITMTAgent {
	agentMap := serv.GetAgentMap()

	agents := make([]agent.ITMTAgent, 0, len(agentMap))
	for _, ag := range agentMap {
		agents = append(agents, ag)
	}

	rand.Shuffle(len(agents), func(i, j int) {
		agents[i], agents[j] = agents[j], agents[i]
	})

	return agents
}

func (serv *GameServer) EstablishInitialObstructions() {
	// Clear agent stored obstructions, then establish current positions
	for _, ag := range serv.GetAgentMap() {
		ag.ClearObstructions()
	}
	for _, ag := range serv.GetAgentMap() {
		ag.BroadcastPosition()
	}
}

func (serv *GameServer) DrainAgents() {
	// Drain agent, and if -ve energy, kill it
	for _, ag := range serv.GetAgentMap() {
		ag.ChangeEnergy(-1)
		if ag.GetEnergy() < 0 {
			serv.RemoveAgent(ag)
		}
	}
}
