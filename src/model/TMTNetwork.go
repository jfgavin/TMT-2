package model

import "github.com/jfgavin/TMT-2/src/config"

type TMTNetwork struct {
	cfg config.SynapseConfig

	inputs map[string]*SynapseInput
	Output *Synapse

	all []*Synapse

	Time float64
}

func (net *TMTNetwork) Step() float64 {

	// 1. Inject input values
	for _, in := range net.inputs {
		in.InjectFromSource()
	}

	// 2. Update membrane potential
	for _, syn := range net.all {
		syn.Advance()
	}

	// 3. Spike propagation
	for _, syn := range net.all {
		syn.Propagate()
	}

	// 4. Advance time, read final output
	net.Time += net.cfg.Dt

	return net.Output.GetOutput()
}

func (net *TMTNetwork) NewSynapse() *Synapse {
	syn := NewSynapse(net.cfg)
	net.all = append(net.all, syn)
	return syn
}

func (net *TMTNetwork) NewInput(name string) *SynapseInput {
	syn := net.NewSynapse()

	input := &SynapseInput{
		Synapse: syn,
		Source:  nil,
	}

	net.inputs[name] = input
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
		inputs: make(map[string]*SynapseInput),
	}

	// Mortality Salience
	elimInput := net.NewInput("eliminations")
	msSynapse := net.NewSynapseBlock(
		[]*Synapse{elimInput.Synapse},
		[]float64{0.8},
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

func (net *TMTNetwork) GetOutput() float64 {
	return net.Output.GetOutput()
}

func (net *TMTNetwork) RegisterInput(name string, function func() float64) {
	net.inputs[name].Source = func() float64 {
		return function()
	}
}
