package infra

import "github.com/jfgavin/TMT-2/src/config"

// === Subtypes ===
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

// === Environment Type ===

type Environment struct {
	cfg  config.EnvironmentConfig
	Grid [][]*Tile
}

func NewEnvironment(cfg config.EnvironmentConfig) *Environment {
	env := &Environment{
		cfg:  cfg,
		Grid: make([][]*Tile, cfg.GridSize),
	}

	for y := range env.Grid {
		env.Grid[y] = make([]*Tile, cfg.GridSize)
	}

	for y := range env.Grid {
		for x := range env.Grid[y] {
			env.Grid[y][x] = NewTile(0)
		}
	}

	return env
}

// === Environment Methods ===

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

func (env *Environment) GetTile(pos Position) (*Tile, bool) {
	tile := env.Grid[pos.Y][pos.X]
	if tile != nil {
		return tile, true
	}

	return nil, false
}

func (env *Environment) BoundPos(pos Position) Position {
	upperBound := env.cfg.GridSize - 1

	if pos.X < 0 {
		pos.X = 0
	} else if pos.X > upperBound {
		pos.X = upperBound
	}

	if pos.Y < 0 {
		pos.Y = 0
	} else if pos.Y > upperBound {
		pos.Y = upperBound
	}

	return pos
}
