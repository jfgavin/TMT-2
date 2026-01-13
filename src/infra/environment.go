package infra

const (
	NumAgents      = 4
	StartingEnergy = 5
	GridSize       = 64
)

type Position struct {
	X, Y int
}

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

	for y := range env.Grid {
		for x := range env.Grid[y] {
			env.Grid[y][x] = NewTile(0)
		}
	}

	return env
}

func (env *Environment) TilePos(tile *Tile) (Position, bool) {
	for y := range env.Grid {
		for x := range env.Grid[y] {
			if env.Grid[y][x] == tile {
				return Position{X: x, Y: y}, true
			}
		}
	}
	return Position{}, false
}
