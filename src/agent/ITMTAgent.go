package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/src/env"
)

type ITMTAgent interface {
	agent.IAgent[ITMTAgent]

	DoMessaging()

	GetName() string

	GetPos() env.Position

	SetPos(pos env.Position)

	ChangeEnergy(energyDelta int)

	GetEnergy() int
}
