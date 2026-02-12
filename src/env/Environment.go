package env

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/config"
)

type Environment struct {
	cfg      config.EnvironmentConfig
	clusters map[uuid.UUID]*Cluster
	graves   map[Position]*Grave
}

func NewEnvironment(cfg config.EnvironmentConfig) *Environment {
	env := &Environment{
		cfg:      cfg,
		clusters: make(map[uuid.UUID]*Cluster),
		graves:   make(map[Position]*Grave),
	}

	env.IntroduceResources()

	return env
}

func (env *Environment) GridSize() int {
	return env.cfg.GridSize
}

func (env *Environment) GetRandPosPadded(padding int) Position {
	return Position{
		X: rand.Intn(env.GridSize()-padding) + padding,
		Y: rand.Intn(env.GridSize()-padding) + padding,
	}
}
