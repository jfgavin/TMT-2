package agent

import (
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

	net.RegisterInput("eliminations", elimCount)
	net.RegisterInput("worldview", tmta.WorldviewScore)

	tmta.tmt = net
}

// Run once per turn
func (tmta *TMTAgent) DriveModel() {

	for i := 0.0; i < 1.0; i += tmta.cfg.Neurons.Dt {
		out := tmta.tmt.Step()
		tmta.ModelOutput = append(tmta.ModelOutput, out)

		if tmta.tmt.Output.DidSpike() {
			tmta.serv.RequestSacrifice(tmta)
		}
	}
}
