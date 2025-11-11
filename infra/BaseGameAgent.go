package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

type BaseGameAgent struct {
	// embed functionality of package
	*agent.BaseAgent[IGameAgent]
	// add additional fields for all agents in simulator
	Pos           Position
	Energy        int
	HasExitedFlag bool
}

// base implementation of DoMessaging
func (bga *BaseGameAgent) DoMessaging() {
	bga.SignalMessagingComplete()
}

func (bga *BaseGameAgent) GetPos() Position {
	return bga.Pos
}

func (bga *BaseGameAgent) GetEnergy() int {
	return bga.Energy
}

func (bga *BaseGameAgent) ResetEnergy() {
	bga.Energy = StartingEnergy
}

func (bga *BaseGameAgent) HasExited() bool {
	return bga.HasExitedFlag
}

func (bga *BaseGameAgent) SetExited(exited bool) {
	bga.HasExitedFlag = exited
}

func (bga *BaseGameAgent) AddEnergy(amount int) {
	bga.Energy += amount
	if bga.Energy > StartingEnergy {
		bga.Energy = StartingEnergy
	}
}

// constructor for BaseGameAgent
func GetBaseGameAgent(funcs agent.IExposedServerFunctions[IGameAgent], pos Position) *BaseGameAgent {
	return &BaseGameAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		Pos:       pos,
		Energy:    StartingEnergy,
	}
}
