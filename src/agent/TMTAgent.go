package agent

import (
	"math"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/src/config"
	"github.com/jfgavin/TMT-2/src/env"
	"github.com/jfgavin/TMT-2/src/model"
)

type TMTAgent struct {
	*agent.BaseAgent[ITMTAgent] `json:"-"`
	cfg                         config.AgentConfig
	env                         env.IEnvironment
	serv                        ServerAPI
	Name                        string
	Pos                         env.Position
	obstructions                map[env.Position]struct{}
	Energy                      int
	tmt                         *model.TMTNetwork
	ModelOutput                 []float64
}

func (tmta *TMTAgent) GetPos() env.Position {
	return tmta.Pos
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

func (tmta *TMTAgent) AssignTMTModel() {
	net := model.NewTMTNetwork(tmta.cfg.Neurons)

	elimCount := func() float64 {
		return float64(tmta.serv.GetEliminationCount())
	}

	net.RegisterInput(model.Eliminations, elimCount)

	tmta.tmt = net
}

func (tmta *TMTAgent) DriveModel() {
	tmta.tmt.Inject()

	steps := int(math.Round(tmta.cfg.Neurons.MsPerStep / tmta.cfg.Neurons.Dt))
	for i := 0; i < steps; i++ {
		out := tmta.tmt.Step()
		tmta.ModelOutput = append(tmta.ModelOutput, out)

		if tmta.tmt.Output.DidSpike() {
			tmta.serv.RequestSacrifice(tmta)
			return
		}
	}
}

func (tmta *TMTAgent) PlayTurn() {
	// Drive the TMT model
	tmta.DriveModel()

	// Always try to harvest resources, then move if that fails
	if !tmta.HarvestResources() {
		tmta.Move()
	}

	// Make sure to signal messaging done after turn
	tmta.SignalMessagingComplete()
}

func NewTMTAgent(funcs agent.IExposedServerFunctions[ITMTAgent], cfg config.AgentConfig, env env.IEnvironment, serv ServerAPI, name string, initPos env.Position) *TMTAgent {
	agent := &TMTAgent{
		BaseAgent:   agent.CreateBaseAgent(funcs),
		cfg:         cfg,
		env:         env,
		serv:        serv,
		Name:        name,
		Pos:         initPos,
		Energy:      cfg.StartingEnergy,
		ModelOutput: make([]float64, 0),
	}

	agent.AssignTMTModel()

	return agent
}
