package server

import (
	"fmt"

	"github.com/jfgavin/TMT-2/src/agent"
	"github.com/jfgavin/TMT-2/src/env"
)

// Functions explicitly callable by the agent via the ServerAPI will be placed here
func (serv *GameServer) GetEliminationCount() int {
	return serv.elims
}

func (serv *GameServer) RequestSacrifice(sacAg agent.ITMTAgent) {
	serv.sacrificeRequests = append(serv.sacrificeRequests, sacAg)
}

func (serv *GameServer) IsObstructed(pos env.Position) bool {
	_, blocked := serv.obstructions[pos]
	return blocked && pos.IsBounded(serv.Env.GridSize())
}

func (serv *GameServer) MoveAgent(ag agent.ITMTAgent, target env.Position) bool {
	agPos := ag.GetPos()
	if agPos.ManhatDist(target) == 1 && !serv.IsObstructed(target) {
		ag.SetPos(target)
		serv.obstructions[target] = struct{}{}
		return true
	}
	return false
}

func (serv *GameServer) RequestChild(parent agent.ITMTAgent) bool {
	if parent.GetEnergy() <= serv.agCfg.ChildCost {
		return false
	}

	pos := parent.GetPos()
	childPos := pos
	for _, adj := range pos.GetShuffledAdjacent() {
		if !serv.IsObstructed(adj) {
			childPos = adj
			break
		}
	}
	if childPos == pos {
		return false
	}

	childName := fmt.Sprintf("%s %d", parent.GetName(), len(parent.GetChildren()))

	child := agent.NewTMTAgent(serv, serv.agCfg, serv.Env, serv, childName, childPos, parent.GetID())
	serv.AddAgent(child)

	parent.ChangeEnergy(-serv.agCfg.ChildCost)
	parent.AddChildID(child.GetID())

	return true
}
