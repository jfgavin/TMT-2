package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type IGameAgent interface {
	// embed functionality from package
	agent.IAgent[IGameAgent]

	DoMessaging()

	GetPos() Position

	SetPos(pos Position)

	GetName() string

	GetEnergy() int

	SetEnergy(energy int)
}
