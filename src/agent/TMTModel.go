package agent

import (
	"github.com/jfgavin/TMT-2/src/model"
)

func (tmta *TMTAgent) WorldviewScore() float64 {
	return 0.0
}

func (tmta *TMTAgent) AssignTMTModel() {
	net := model.NewTMTNetwork(tmta.cfg.Synapses)

	net.Inputs["eliminations"].Source = func() float64 {
		return float64(tmta.serv.GetEliminationCount() * 100)
	}

	net.Inputs["worldview"].Source = func() float64 {
		return tmta.WorldviewScore()
	}

	tmta.tmt = net
}

// Run once per turn
func (tmta *TMTAgent) DriveModel() {

	for i := 0.0; i < 1.0; i += tmta.cfg.Synapses.Dt {
		tmta.tmt.Step()
		tmta.ModelOutput = append(tmta.ModelOutput, tmta.tmt.GetOutput())
	}
}
