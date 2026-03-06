package server

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
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

func (serv *GameServer) UpdateObstructions() {
	// Clear stored obstructions
	serv.obstructions = make(map[env.Position]struct{})

	// Re-establish current obstructions
	for _, ag := range serv.GetAgentMap() {
		serv.obstructions[ag.GetPos()] = struct{}{}
	}
}

func (serv *GameServer) HandleAgentMortality() {
	serv.elims = 0
	for _, ag := range serv.GetAgentMap() {
		// Drain agent, then kill if no energy
		ag.ChangeEnergy(-1)
		if ag.GetEnergy() < 0 {
			serv.KillAgent(ag)
		}
	}

	// Sacrifice agents which requested it
	for _, ag := range serv.sacrificeRequests {
		serv.SacrificeAgent(ag)
	}

	// Clear sacrifice buffer
	serv.sacrificeRequests = nil
}

func (serv *GameServer) KillAgent(ag agent.ITMTAgent) {
	serv.Env.PlaceTombstone(ag.GetPos())
	serv.RemoveAgent(ag)
	serv.elims++
}

func (serv *GameServer) SacrificeAgent(ag agent.ITMTAgent) {
	serv.Env.PlaceMemorial(ag.GetPos())
	serv.RemoveAgent(ag)
	serv.elims++
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
		ga := agent.NewTMTAgent(serv, serv.agCfg, serv.Env, serv, fmt.Sprintf("Agent %d", i), pos, uuid.Nil)
		serv.AddAgent(ga)
	}
}
