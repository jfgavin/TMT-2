package env

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/config"
)

type Environment struct {
	cfg      config.EnvironmentConfig
	Grid     [][]*Tile
	clusters map[uuid.UUID]*Cluster
}

func NewEnvironment(cfg config.EnvironmentConfig) *Environment {
	env := &Environment{
		cfg:      cfg,
		Grid:     make([][]*Tile, cfg.GridSize),
		clusters: make(map[uuid.UUID]*Cluster),
	}

	for y := range env.Grid {
		env.Grid[y] = make([]*Tile, cfg.GridSize)
		for x := range env.Grid[y] {
			env.Grid[y][x] = NewTile()
		}
	}

	env.IntroduceResources()

	return env
}

func (env *Environment) GridSize() int {
	return env.cfg.GridSize
}

func (env *Environment) GetTile(pos Position) (*Tile, bool) {
	if pos.Y < 0 || pos.Y >= len(env.Grid) {
		return nil, false
	}
	row := env.Grid[pos.Y]

	if pos.X < 0 || pos.X >= len(row) {
		return nil, false
	}

	return row[pos.X], true
}

func (env *Environment) GetRandPosPadded(padding int) Position {
	return Position{
		X: rand.Intn(env.GridSize()-padding) + padding,
		Y: rand.Intn(env.GridSize()-padding) + padding,
	}
}
