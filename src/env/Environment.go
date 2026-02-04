package env

import (
	"math"
	"math/rand"

	"github.com/jfgavin/TMT-2/src/config"
)

type Environment struct {
	cfg  config.EnvironmentConfig
	Grid [][]Tile
}

func NewEnvironment(cfg config.EnvironmentConfig) *Environment {
	env := &Environment{
		cfg:  cfg,
		Grid: make([][]Tile, cfg.GridSize),
	}

	for y := range env.Grid {
		env.Grid[y] = make([]Tile, cfg.GridSize)
		for x := range env.Grid[y] {
			env.Grid[y][x] = NewTile(0)
		}
	}

	return env
}

func (env *Environment) GridSize() int {
	return env.cfg.GridSize
}

func (env *Environment) GetTile(pos Position) (Tile, bool) {
	if pos.Y < 0 || pos.Y >= len(env.Grid) {
		return Tile{}, false
	}
	row := env.Grid[pos.Y]

	if pos.X < 0 || pos.X >= len(row) {
		return Tile{}, false
	}

	return row[pos.X], true
}

func (env *Environment) ChangeResources(pos Position, delta int) bool {
	if pos.IsBounded(env.GridSize()) {
		env.Grid[pos.Y][pos.X].ChangeResources(delta)
		return true
	}
	return false
}

func (env *Environment) GetRandPosPadded(padding int) Position {
	return Position{
		X: rand.Intn(env.GridSize()-padding) + padding,
		Y: rand.Intn(env.GridSize()-padding) + padding,
	}
}

func (env *Environment) IntroduceResources() {
	cfg := env.cfg.Resources
	radius, lambda := float64(cfg.Radius), float64(cfg.Lambda)

	maxTerm := 1 - math.Exp(-radius/lambda)

	// Find cluster centres
	// Ensure their most extreme points are not outside of the grid
	centres := make([]Position, cfg.ClusterCount)
	for i := 0; i < cfg.ClusterCount; i++ {
		centres[i] = env.GetRandPosPadded(cfg.Radius)
	}

	// Randomly place resources, one-by-one
	for cfg.ResourceCount > 0 {
		chosenCentre := centres[rand.Intn(cfg.ClusterCount)]

		// Random angle
		theta := rand.Float64() * 2 * math.Pi

		// Random distance from centre
		u := rand.Float64()
		dist := -lambda * math.Log(1-u*maxTerm)

		// Clamping
		if dist > radius {
			dist = radius
		}

		// Final position of new resource
		x := float64(chosenCentre.X) + dist*math.Cos(theta)
		y := float64(chosenCentre.Y) + dist*math.Sin(theta)

		newPos := Position{
			X: int(math.Round(x)),
			Y: int(math.Round(y)),
		}
		newPos.Bound(env.GridSize())

		// Modify tile
		if ok := env.ChangeResources(newPos, 1); ok {
			cfg.ResourceCount--
		}
	}
}
