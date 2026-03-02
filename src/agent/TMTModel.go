package agent

import "github.com/jfgavin/TMT-2/src/model"

func (tmta *TMTAgent) WorldviewScore() float64 {
	return 0.5
}

func (tmta *TMTAgent) AssignTMTModel() {
	net := model.NewTMTNetwork(tmta.cfg.Synapses)

	net.Inputs["eliminations"].Source = func() float64 {
		return float64(tmta.serv.GetEliminationCount())
	}

	net.Inputs["worldview"].Source = func() float64 {
		return tmta.WorldviewScore()
	}

	tmta.tmt = net
}

func (tmta *TMTAgent) DriveModel() {
	for i := 0; i < 10; i++ {
		tmta.tmt.Step()
	}
}
