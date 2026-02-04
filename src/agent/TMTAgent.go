package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/src/config"
	"github.com/jfgavin/TMT-2/src/env"
)

type TMTAgent struct {
	*agent.BaseAgent[ITMTAgent] `json:"-"`
	cfg                         config.AgentConfig
	env                         env.IEnvironment
	Name                        string
	Pos                         env.Position
	obstructions                map[env.Position]struct{}
	Energy                      int
}

func (tmta *TMTAgent) ChangeEnergy(energyDelta int) {
	tmta.Energy += energyDelta
}

func (tmta *TMTAgent) GetEnergy() int {
	return tmta.Energy
}

func (tmta *TMTAgent) ClearObstructions() {
	tmta.obstructions = make(map[env.Position]struct{})
}

func (tmta *TMTAgent) BroadcastPosition() {
	msg := tmta.NewObstructionMessage(tmta.Pos)
	tmta.BroadcastSynchronousMessage(msg)
}

func (tmta *TMTAgent) PlayTurn() {
	if !tmta.HarvestResources() {
		tmta.TargetedMove()
	}

	// Make sure to signal messaging done after turn
	tmta.SignalMessagingComplete()
}

func NewTMTAgent(funcs agent.IExposedServerFunctions[ITMTAgent], cfg config.AgentConfig, env env.IEnvironment, name string, initPos env.Position) *TMTAgent {
	return &TMTAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		cfg:       cfg,
		env:       env,
		Name:      name,
		Pos:       initPos,
		Energy:    cfg.StartingEnergy,
	}
}
