package model

import "github.com/jfgavin/TMT-2/src/config"

type TMTNetwork struct {
	cfg config.SynapseConfig

	Inputs map[string]*SynapseInput
	Output *Synapse

	All []*Synapse

	Time float64
}

func (net *TMTNetwork) Step() float64 {

	// 1. Inject input values
	for _, in := range net.Inputs {
		in.InjectFromSource()
	}

	// 2. Advance all synapses
	for _, syn := range net.All {
		syn.Advance()
	}

	// 3. Propagate outputs
	for _, syn := range net.All {
		out := syn.GetOutput()
		for _, conn := range syn.Outgoing {
			conn.Target.Inject(conn.Weight * out)
		}
	}

	// 4. Advance time, read final output
	net.Time += net.cfg.Dt

	return net.Output.GetOutput()
}

func (net *TMTNetwork) NewSynapse() *Synapse {
	syn := NewSynapse(net.cfg)
	net.All = append(net.All, syn)
	return syn
}

func (net *TMTNetwork) NewInput(name string) *SynapseInput {
	syn := net.NewSynapse()

	input := &SynapseInput{
		Synapse: syn,
		Source:  nil,
	}

	net.Inputs[name] = input
	return input
}

// Used when multiple synapses target one synapse
func (net *TMTNetwork) NewSynapseBlock(inputs []*Synapse, weights []float64) *Synapse {
	out := net.NewSynapse()
	for i, syn := range inputs {
		syn.AddConnection(out, weights[i])
	}

	return out
}

// A method of taking function handles as network inputs
type SynapseInput struct {
	*Synapse
	Source func() float64
}

func (sin *SynapseInput) InjectFromSource() {
	if sin.Source != nil {
		sin.Inject(sin.Source())
	}
}

func NewTMTNetwork(cfg config.SynapseConfig) *TMTNetwork {
	net := &TMTNetwork{
		cfg:    cfg,
		Inputs: make(map[string]*SynapseInput),
	}

	// Mortality Salience
	elimInput := net.NewInput("eliminations")
	msSynapse := net.NewSynapseBlock(
		[]*Synapse{elimInput.Synapse},
		[]float64{0.1},
	)

	// Worldview
	worldviewValidationInput := net.NewInput("worldview")
	wvSynapse := net.NewSynapseBlock(
		[]*Synapse{worldviewValidationInput.Synapse},
		[]float64{0.5},
	)

	// Output
	out := net.NewSynapseBlock(
		[]*Synapse{
			msSynapse,
			wvSynapse,
		},
		[]float64{
			0.5,
			0.5,
		},
	)

	net.Output = out
	return net
}
