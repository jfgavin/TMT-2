package agent

import "fmt"

func (tmta *TMTAgent) DriveModel() {
	elims := tmta.serv.GetEliminationCount()
	if elims > 0 {
		elims_flt := float64(elims)
		tmta.Syn.Inject(elims_flt)
		fmt.Printf("Injected an elim spike for agent: %s\n", tmta.Name)
	}

	for i := 0; i < 100; i++ {
		tmta.Syn.Advance()

	}

	if tmta.Name == "Agent 0" && tmta.Energy <= 1 {
		fmt.Println("AGENT 0 HIT")
		tmta.Syn.PlotSynapse()
	}

}
