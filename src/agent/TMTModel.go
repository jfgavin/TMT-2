package agent

import (
	"fmt"

	"github.com/jfgavin/TMT-2/src/model"
)

func (tmta *TMTAgent) WorldviewScore() float64 {
	return 0.0
}

func (tmta *TMTAgent) AssignTMTModel() {
	net := model.NewTMTNetwork(tmta.cfg.Synapses)

	elimCount := func() float64 {
		return float64(tmta.serv.GetEliminationCount())
	}

	net.RegisterInput("eliminations", elimCount)
	net.RegisterInput("worldview", tmta.WorldviewScore)

	tmta.tmt = net
}

// Run once per turn
func (tmta *TMTAgent) DriveModel() {

	for i := 0.0; i < 1.0; i += tmta.cfg.Synapses.Dt {
		out := tmta.tmt.Step()
		tmta.ModelOutput = append(tmta.ModelOutput, out)

		if out >= 1.0 {
			tmta.serv.RequestSacrifice(tmta)
			fmt.Printf("%s has fired!!!\n", tmta.Name)
		}
	}
}
