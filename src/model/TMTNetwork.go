package model

import "github.com/jfgavin/TMT-2/src/config"

type TMTInputSource int

const (
	Eliminations TMTInputSource = iota
	WorldviewValidation
)

type TMTNetwork struct {
	cfg config.NeuronConfig

	inputs map[TMTInputSource]*NeuronInput
	Output *Neuron

	all []*Neuron

	Time float64
}

func (net *TMTNetwork) Inject() {
	// 1. Inject input values
	for _, in := range net.inputs {
		in.InjectFromSource()
	}
}

func (net *TMTNetwork) Step() float64 {
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

func (net *TMTNetwork) NewNeuron() *Neuron {
	syn := NewNeuron(net.cfg)
	net.all = append(net.all, syn)
	return syn
}

func (net *TMTNetwork) NewInput(source TMTInputSource) *NeuronInput {
	syn := net.NewNeuron()

	input := &NeuronInput{
		Neuron: syn,
		Source: nil,
	}

	net.inputs[source] = input
	return input
}

// Used when multiple Neurons target one Neuron
func (net *TMTNetwork) NewNeuronBlock(inputs []*Neuron, weights []float64) *Neuron {
	out := net.NewNeuron()
	for i, syn := range inputs {
		syn.AddConnection(out, weights[i])
	}

	return out
}

// A way of taking function handles as network inputs
type NeuronInput struct {
	*Neuron
	Source func() float64
}

func (sin *NeuronInput) InjectFromSource() {
	if sin.Source != nil {
		sin.Inject(sin.Source())
	}
}

func NewTMTNetwork(cfg config.NeuronConfig) *TMTNetwork {
	net := &TMTNetwork{
		cfg:    cfg,
		inputs: make(map[TMTInputSource]*NeuronInput),
	}

	// Mortality Salience
	elimInput := net.NewInput(Eliminations)
	msNeuron := net.NewNeuronBlock(
		[]*Neuron{elimInput.Neuron},
		[]float64{1.0},
	)

	// Worldview
	worldviewValidationInput := net.NewInput(WorldviewValidation)
	wvNeuron := net.NewNeuronBlock(
		[]*Neuron{worldviewValidationInput.Neuron},
		[]float64{0.5},
	)

	// Output
	out := net.NewNeuronBlock(
		[]*Neuron{
			msNeuron,
			wvNeuron,
		},
		[]float64{
			1.0,
			0.5,
		},
	)

	net.Output = out
	return net
}

func (net *TMTNetwork) GetOutput() float64 {
	return net.Output.GetOutput()
}

func (net *TMTNetwork) RegisterInput(source TMTInputSource, function func() float64) {
	net.inputs[source].Source = func() float64 {
		return function()
	}
}
