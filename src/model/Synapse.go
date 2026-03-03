package model

import (
	"math"

	"github.com/jfgavin/TMT-2/src/config"
)

type Synapse struct {
	TauRise  float64
	TauDecay float64
	Dt       float64

	Threshold float64

	V float64 // membrane potential

	riseState  float64 // rise  accumulator
	decayState float64 // decay accumulator
	riseRate   float64 // exp(-Dt/TauRise)
	decayRate  float64 // exp(-Dt/TauDecay)
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

func NewSynapse(cfg config.SynapseConfig) *Synapse {
	riseRate := math.Exp(-cfg.Dt / cfg.TauRise)
	decayRate := math.Exp(-cfg.Dt / cfg.TauDecay)

	// Analytical peak time and peak value of (exp(-t/τd) - exp(-t/τr))
	// t_peak = (τr*τd)/(τd-τr) * ln(τd/τr)
	// peak   = exp(-t_peak/τd) - exp(-t_peak/τr)
	tPeak := (cfg.TauRise * cfg.TauDecay) / (cfg.TauDecay - cfg.TauRise) *
		math.Log(cfg.TauDecay/cfg.TauRise)
	peak := math.Exp(-tPeak/cfg.TauDecay) - math.Exp(-tPeak/cfg.TauRise)
	normFactor := 1.0 / peak // multiply injected weight by this

	return &Synapse{
		TauRise:    cfg.TauRise,
		TauDecay:   cfg.TauDecay,
		Dt:         cfg.Dt,
		Threshold:  1.0,
		riseRate:   riseRate,
		decayRate:  decayRate,
		normFactor: normFactor,
		Outgoing:   make([]Connection, 0),
	}
}

func (syn *Synapse) Inject(weight float64) {
	scaled := weight * syn.normFactor
	syn.riseState += scaled
	syn.decayState += scaled
}

func (syn *Synapse) Advance() {
	syn.riseState *= syn.riseRate
	syn.decayState *= syn.decayRate
	syn.V += syn.decayState - syn.riseState
}

func (syn *Synapse) Propagate() {
	syn.spikedFlag = false
	if syn.V >= syn.Threshold {
		for _, conn := range syn.Outgoing {
			conn.Target.Inject(conn.Weight)
		}
		syn.V = 0
		syn.spikedFlag = true
	}
}

func (syn *Synapse) GetOutput() float64 {
	return syn.V
}

func (syn *Synapse) DidSpike() bool {
	return syn.spikedFlag
}

func (syn *Synapse) AddConnection(target *Synapse, weight float64) {
	syn.Outgoing = append(syn.Outgoing, Connection{Weight: weight, Target: target})
}
