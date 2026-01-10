package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

type GameAgent struct {
	*agent.BaseAgent[IGameAgent]
	Name   string
	Pos    Position
	Energy int
}

func (bga *GameAgent) DoMessaging() {
	bga.SignalMessagingComplete()
}

func (bga *GameAgent) GetPos() Position {
	return bga.Pos
}

func (bga *GameAgent) SetPos(pos Position) {
	bga.Pos = pos
}

func (bga *GameAgent) GetEnergy() int {
	return bga.Energy
}

func (bga *GameAgent) SetEnergy(energy int) {
	bga.Energy = energy
}

func NewGameAgent(funcs agent.IExposedServerFunctions[IGameAgent], name string, pos Position) *GameAgent {
	return &GameAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		Name:      name,
		Pos:       pos,
		Energy:    StartingEnergy,
	}
}
