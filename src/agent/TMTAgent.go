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
	obstructions                []env.Position
	Energy                      int
	percetiveRange              int
}

func (tmta *TMTAgent) ChangeEnergy(energyDelta int) {
	tmta.Energy += energyDelta
}

func (tmta *TMTAgent) GetEnergy() int {
	return tmta.Energy
}

func (tmta *TMTAgent) ClearObstructions() {
	tmta.obstructions = nil
}

func (tmta *TMTAgent) BroadcastPosition() {
	msg := tmta.NewObstructionMessage(tmta.Pos)
	tmta.BroadcastSynchronousMessage(msg)
}

func (tmta *TMTAgent) PlayTurn() {
	currTile, found := tmta.env.GetTile(tmta.Pos)
	if !found {
		// Agent is not on grid. Just exit
		return
	}

	nextStep := tmta.getUnobstructedBestStep()

	if currTile.GetResources() > 0 {
		tmta.HarvestResources()
	} else {
		tmta.Pos = nextStep
		msg := tmta.NewObstructionMessage(nextStep)
		tmta.BroadcastSynchronousMessage(msg)
	}

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
