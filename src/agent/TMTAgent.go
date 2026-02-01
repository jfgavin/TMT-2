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
	Energy                      int
	percetiveRange              int
}

func (tmta *TMTAgent) DoMessaging() {
	tmta.SignalMessagingComplete()
}

func (tmta *TMTAgent) GetName() string {
	return tmta.Name
}

func (tmta *TMTAgent) ChangeEnergy(energyDelta int) {
	tmta.Energy += energyDelta
}

func (tmta *TMTAgent) GetEnergy() int {
	return tmta.Energy
}

func NewTMTAgent(funcs agent.IExposedServerFunctions[ITMTAgent], cfg config.AgentConfig, env env.IEnvironment, name string, initPos env.Position) *TMTAgent {
	return &TMTAgent{
		BaseAgent:      agent.CreateBaseAgent(funcs),
		cfg:            cfg,
		env:            env,
		Name:           name,
		Pos:            initPos,
		Energy:         cfg.StartingEnergy,
		percetiveRange: cfg.PerceptiveRange,
	}
}
