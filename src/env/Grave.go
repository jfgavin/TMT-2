package env

type Grave struct {
	AgeRemaining int
}

func NewGrave(ageRem int) *Grave {
	return &Grave{
		AgeRemaining: ageRem,
	}
}

func (grv *Grave) AgeRem() int {
	return grv.AgeRemaining
}

func (grv *Grave) Tick() {
	grv.AgeRemaining--
}

func (env *Environment) GetGraves() map[Position]*Grave {
	return env.graves
}

func (env *Environment) PlaceGrave(pos Position) {
	env.graves[pos] = NewGrave(env.cfg.GraveLifespan)
}

func (env *Environment) TickGraves() {
	for pos, grave := range env.graves {
		grave.Tick()
		if grave.AgeRem() <= 0 {
			delete(env.graves, pos)
		}
	}
}
