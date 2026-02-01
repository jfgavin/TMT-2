package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

type ITMTAgent interface {
	agent.IAgent[ITMTAgent]

	DoMessaging()

	GetName() string

	Move()

	ChangeEnergy(energyDelta int)

	GetEnergy() int
}
