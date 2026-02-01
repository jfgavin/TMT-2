package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

type ITMTAgent interface {
	agent.IAgent[ITMTAgent]

	DoMessaging()

	PlayTurn()

	ChangeEnergy(energyDelta int)

	GetEnergy() int
}
