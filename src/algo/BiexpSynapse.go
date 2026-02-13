package algo

import "math"

type BiexpSynapse struct {
	TauRise  float64
	TauDecay float64
	Dt       float64

	riseState  float64
	decayState float64

	riseFactor  float64
	decayFactor float64

	Output   float64
	Outgoing []Connection
}

type Connection struct {
	Weight float64
	Target SynapticTarget
}

type SynapticTarget interface {
	Inject(weight float64)
}

func NewBiexpSynapse(tauRise, tauDecay, dt float64) *BiexpSynapse {
	return &BiexpSynapse{
		TauRise:     tauRise,
		TauDecay:    tauDecay,
		Dt:          dt,
		riseFactor:  math.Exp(-dt / tauRise),
		decayFactor: math.Exp(-dt / tauDecay),
	}
}

func (syn *BiexpSynapse) Inject(w float64) {
	syn.riseState += w
	syn.decayState += w
}

func (syn *BiexpSynapse) Advance() float64 {
	syn.riseState *= syn.riseFactor
	syn.decayState *= syn.decayFactor
	return syn.decayState - syn.riseState
}

func (syn *BiexpSynapse) Propagate(output float64) {
	for _, conn := range syn.Outgoing {
		conn.Target.Inject(conn.Weight * syn.Output)
	}
}
