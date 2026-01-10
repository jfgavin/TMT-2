package infra

const (
	NumAgents      = 4
	StartingEnergy = 5
	GridSize       = 8
)

type Tile struct {
	Resources int
}

func NewTile(resources int) *Tile {
	return &Tile{
		Resources: resources,
	}
}

type Environment struct {
	Grid [][]*Tile
}

func NewEnvironment() *Environment {
	env := &Environment{
		Grid: make([][]*Tile, GridSize),
	}

	for y := range env.Grid {
		env.Grid[y] = make([]*Tile, GridSize)
	}

	return env
}
