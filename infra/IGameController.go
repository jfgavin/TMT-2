package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type IGameAgent interface {
	// embed functionality from package
	agent.IAgent[IGameAgent]
	// perform messaging action
	DoMessaging()
	// get 2D grid pos
	GetPos() Position

	SetAction(action Action)

	GetAction() Action

	SetMovementDir(movement Direction)

	GetMovementDir() Direction

	DoMovement()

	GetName() string

	GetEnergy() int

	ResetEnergy()

	// perform action
	DoAction()

	// check if hero has exited
	HasExited() bool

	// check if hero has loot
	HasLoot() bool

	// set exit status
	SetExited(bool)

	// add energy (for loot shop interactions)
	AddEnergy(int)
}
