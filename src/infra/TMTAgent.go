package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/src/config"
)

type TMTAgent struct {
	*agent.BaseAgent[ITMTAgent] `json:"-"`
	cfg                         config.AgentConfig
	env                         IEnvironment
	Name                        string
	Pos                         Position
	Energy                      int
}

func (tmta *TMTAgent) DoMessaging() {
	tmta.SignalMessagingComplete()
}

func (tmta *TMTAgent) GetName() string {
	return tmta.Name
}

func (tmta *TMTAgent) GetPos() Position {
	return tmta.Pos
}

func (tmta *TMTAgent) SetPos(pos Position) {
	tmta.Pos = tmta.env.BoundPos(pos)
}

func (tmta *TMTAgent) GetEnergy() int {
	return tmta.Energy
}

func (tmta *TMTAgent) SetEnergy(energy int) {
	tmta.Energy = energy
}

func (tmta *TMTAgent) ResetEnergy() {
	tmta.Energy = tmta.cfg.StartingEnergy
}

func NewTMTAgent(funcs agent.IExposedServerFunctions[ITMTAgent], cfg config.AgentConfig, env IEnvironment, name string, initPos Position) *TMTAgent {
	return &TMTAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		cfg:       cfg,
		env:       env,
		Name:      name,
		Pos:       initPos,
		Energy:    cfg.StartingEnergy,
	}
}
