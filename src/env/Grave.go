package env

type GraveType int

const (
	Tombstone GraveType = iota
	Memorial
)

type Grave struct {
	AgeRemaining int
	Type         GraveType
}

func NewGrave(ageRem int, graveType GraveType) *Grave {
	return &Grave{
		AgeRemaining: ageRem,
		Type:         graveType,
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

func (env *Environment) PlaceTombstone(pos Position) {
	env.graves[pos] = NewGrave(
		env.cfg.GraveLifespan,
		Tombstone,
	)
}

func (env *Environment) PlaceMemorial(pos Position) {
	env.graves[pos] = NewGrave(
		env.cfg.GraveLifespan,
		Memorial,
	)
}

func (env *Environment) TickGraves() {
	for pos, grave := range env.graves {
		grave.Tick()
		if grave.AgeRem() <= 0 {
			delete(env.graves, pos)
		}
	}
}
