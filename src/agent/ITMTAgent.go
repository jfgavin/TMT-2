package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/src/env"
)

type ITMTAgent interface {
	agent.IAgent[ITMTAgent]

	PlayTurn()

	GetPos() env.Position

	SetPos(env.Position)

	ChangeEnergy(energyDelta int)

	GetEnergy() int
}
