package model

import (
	"math"

	"github.com/jfgavin/TMT-2/src/config"
)

type Synapse struct {
	Tau       float64
	Dt        float64
	Threshold float64

	V           float64 // membrane potential
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
		Tau:         cfg.TauDecay,
		Dt:          cfg.Dt,
		Threshold:   1.0,
		decayFactor: math.Exp(-cfg.Dt / cfg.TauDecay),
		Outgoing:    make([]Connection, 0),
	}
}

func (syn *Synapse) Inject(w float64) {
	syn.V += w
}

func (syn *Synapse) Advance() {
	syn.V *= syn.decayFactor
}

func (syn *Synapse) GetOutput() float64 {
	return syn.V
}

func (syn *Synapse) Propagate() {
	if syn.V >= syn.Threshold {

		// Emit spike (value 1.0)
		for _, conn := range syn.Outgoing {
			conn.Target.Inject(conn.Weight)
		}

		// Reset membrane
		syn.V = 0
	}
}
