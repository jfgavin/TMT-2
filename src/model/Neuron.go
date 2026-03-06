package model

import (
	"math"

	"github.com/jfgavin/TMT-2/src/config"
)

type Neuron struct {
	TauRise  float64
	TauDecay float64
	Dt       float64

	Threshold float64

	V float64 // membrane potential

	riseState  float64 // rise accumulator
	decayState float64 // decay accumulator
	riseRate   float64 // exp(-Dt/TauRise)
	decayRate  float64 // exp(-Dt/TauDecay)
	membRate   float64
	normFactor float64 // peak-normalisation so max(g)≈1 per unit weight

	spikedFlag bool

	Outgoing []Connection
}
type Connection struct {
	Weight float64
	Target SynapticTarget
}

type SynapticTarget interface {
	Inject(weight float64)
}

func NewNeuron(cfg config.NeuronConfig) *Neuron {
	if cfg.TauRise == cfg.TauDecay {
		panic("TauRise == TauDecay, creating div by 0 bug. The bi-exponential kernel requires TauRise < TauDecay")
	} else if cfg.TauRise > cfg.TauDecay {
		panic("TauRise > TauDecay. The bi-exponential kernel requires TauRise < TauDecay")
	}

	riseRate := math.Exp(-cfg.Dt / cfg.TauRise)
	decayRate := math.Exp(-cfg.Dt / cfg.TauDecay)
	tPeak := (cfg.TauRise * cfg.TauDecay) / (cfg.TauDecay - cfg.TauRise) *
		math.Log(cfg.TauDecay/cfg.TauRise)
	peak := math.Exp(-tPeak/cfg.TauDecay) - math.Exp(-tPeak/cfg.TauRise)
	normFactor := 1.0 / peak

	return &Neuron{
		TauRise:    cfg.TauRise,
		TauDecay:   cfg.TauDecay,
		Dt:         cfg.Dt,
		Threshold:  1.0,
		riseRate:   riseRate,
		decayRate:  decayRate,
		membRate:   math.Exp(-cfg.Dt / cfg.TauMemb),
		normFactor: normFactor,
		Outgoing:   make([]Connection, 0),
	}
}

func (syn *Neuron) Inject(weight float64) {
	scaled := weight * syn.normFactor
	syn.riseState += scaled
	syn.decayState += scaled
}

func (syn *Neuron) Advance() {
	syn.riseState *= syn.riseRate
	syn.decayState *= syn.decayRate
	syn.V = (syn.V * syn.membRate) + (syn.decayState - syn.riseState)
}

func (syn *Neuron) Propagate() {
	syn.spikedFlag = false
	if syn.V >= syn.Threshold {
		for _, conn := range syn.Outgoing {
			conn.Target.Inject(conn.Weight)
		}
		syn.V = 0
		syn.riseState = 0
		syn.decayState = 0
		syn.spikedFlag = true
	}
}

func (syn *Neuron) GetOutput() float64 {
	return syn.V
}

func (syn *Neuron) DidSpike() bool {
	return syn.spikedFlag
}

func (syn *Neuron) AddConnection(target *Neuron, weight float64) {
	syn.Outgoing = append(syn.Outgoing, Connection{Weight: weight, Target: target})
}
