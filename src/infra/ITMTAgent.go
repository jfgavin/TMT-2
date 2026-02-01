package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type ITMTAgent interface {
	agent.IAgent[ITMTAgent]

	DoMessaging()

	GetName() string

	GetPos() Position

	SetPos(pos Position)

	ChangeEnergy(energyDelta int)

	GetEnergy() int
}
