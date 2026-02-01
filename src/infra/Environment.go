package infra

import (
	"math"
	"math/rand"

	"github.com/jfgavin/TMT-2/src/config"
)

// === Subtypes ===
type Position struct {
	X, Y int
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

func (env *Environment) GetRandPos() Position {
	return Position{
		X: rand.Intn(env.cfg.GridSize),
		Y: rand.Intn(env.cfg.GridSize),
	}
}

func (env *Environment) GetRandPosPadded(padding int) Position {
	return Position{
		X: rand.Intn(env.cfg.GridSize-padding) + padding,
		Y: rand.Intn(env.cfg.GridSize-padding) + padding,
	}
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

		newPos := env.BoundPos(Position{
			X: int(math.Round(x)),
			Y: int(math.Round(y)),
		})

		// Modify tile
		tile, found := env.GetTile(newPos)

		if found {
			tile.AddResources(1)
			cfg.ResourceCount--
		}
	}
}
