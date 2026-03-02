package model

import (
	"math"

	"github.com/jfgavin/TMT-2/src/config"
)

type Synapse struct {
	TauRise  float64
	TauDecay float64
	Dt       float64

	riseState  float64
	decayState float64

	riseFactor  float64
	decayFactor float64

	Outgoing []Connection
}

type Connection struct {
	Weight float64
	Target SynapticTarget
}

func (syn *Synapse) AddConnection(target *Synapse, weight float64) {
	syn.Outgoing = append(syn.Outgoing, Connection{Weight: weight, Target: target})
}

type SynapticTarget interface {
	Inject(weight float64)
}

func NewSynapse(cfg config.SynapseConfig) *Synapse {
	return &Synapse{
		TauRise:     cfg.TauRise,
		TauDecay:    cfg.TauDecay,
		Dt:          cfg.Dt,
		riseFactor:  math.Exp(-cfg.Dt / cfg.TauRise),
		decayFactor: math.Exp(-cfg.Dt / cfg.TauDecay),
		Outgoing:    make([]Connection, 0),
	}
}

func (syn *Synapse) Inject(w float64) {
	syn.riseState += w
	syn.decayState += w
}

func (syn *Synapse) Advance() {
	syn.riseState *= syn.riseFactor
	syn.decayState *= syn.decayFactor
}

func (syn *Synapse) GetOutput() float64 {
	return syn.decayState - syn.riseState
}

func (syn *Synapse) Propagate(output float64) {
	for _, conn := range syn.Outgoing {
		conn.Target.Inject(conn.Weight * syn.GetOutput())
	}

}
