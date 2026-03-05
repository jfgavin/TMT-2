package agent

import (
	"math"

	"github.com/jfgavin/TMT-2/src/model"
)

func (tmta *TMTAgent) WorldviewScore() float64 {
	return 0.0
}

func (tmta *TMTAgent) AssignTMTModel() {
	net := model.NewTMTNetwork(tmta.cfg.Neurons)

	elimCount := func() float64 {
		return float64(tmta.serv.GetEliminationCount())
	}

	net.RegisterInput(model.Eliminations, elimCount)
	net.RegisterInput(model.WorldviewValidation, tmta.WorldviewScore)

	tmta.tmt = net
}

// Run once per turn
func (tmta *TMTAgent) DriveModel() {
	tmta.tmt.Inject()

	steps := int(math.Round(tmta.cfg.Neurons.MsPerStep / tmta.cfg.Neurons.Dt))
	for i := 0; i < steps; i++ {
		out := tmta.tmt.Step()
		tmta.ModelOutput = append(tmta.ModelOutput, out)

		if tmta.tmt.Output.DidSpike() {
			// Hitting this is currently failing to transmit the modification to the agent since the last game state
			// Request should store agent in a server buffer until the next turn, where the kill is executed
			tmta.serv.RequestSacrifice(tmta)
			return
		}
	}
}
