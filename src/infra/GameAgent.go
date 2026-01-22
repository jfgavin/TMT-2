package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

type GameAgent struct {
	*agent.BaseAgent[IGameAgent]
	Name   string
	Tile   *Tile
	Energy int
}

func (bga *GameAgent) DoMessaging() {
	bga.SignalMessagingComplete()
}

func (bga *GameAgent) GetTile() *Tile {
	return bga.Tile
}

func (bga *GameAgent) SetTile(tile *Tile) {
	bga.Tile = tile
}

func (bga *GameAgent) GetEnergy() int {
	return bga.Energy
}

func (bga *GameAgent) SetEnergy(energy int) {
	bga.Energy = energy
}

func NewGameAgent(funcs agent.IExposedServerFunctions[IGameAgent], name string, tile *Tile, energy int) *GameAgent {
	return &GameAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		Name:      name,
		Tile:      tile,
		Energy:    energy,
	}
}
