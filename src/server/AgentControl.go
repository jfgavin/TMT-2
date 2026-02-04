package server

import (
	"fmt"
	"math/rand"

	"github.com/jfgavin/TMT-2/src/agent"
	"github.com/jfgavin/TMT-2/src/env"
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

func (serv *GameServer) IntroduceAgents() {
	gs := serv.Env.GridSize()
	// Make 1D list of positions
	positions := make([]env.Position, 0, gs*gs)
	for x := 0; x < gs; x++ {
		for y := 0; y < gs; y++ {
			positions = append(positions, env.Position{X: x, Y: y})
		}
	}

	// Shuffle positions
	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	// Introduce agents, assigning random unique positions
	for i := 0; i < serv.cfg.NumAgents; i++ {
		pos := positions[i]
		ga := agent.NewTMTAgent(serv, serv.agCfg, serv.Env, fmt.Sprintf("Agent %d", i), pos)
		serv.AddAgent(ga)
	}
}
