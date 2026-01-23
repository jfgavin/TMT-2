package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type ITMTAgent interface {
	// embed functionality from package
	agent.IAgent[ITMTAgent]

	DoMessaging()

	GetName() string

	GetPos() Position

	SetPos(pos Position)

	GetEnergy() int

	SetEnergy(energy int)

	ResetEnergy()
}
