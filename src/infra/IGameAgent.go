package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type IGameAgent interface {
	// embed functionality from package
	agent.IAgent[IGameAgent]

	DoMessaging()

	GetName() string

	GetPos() Position

	SetPos(pos Position)

	GetEnergy() int

	SetEnergy(energy int)
}
