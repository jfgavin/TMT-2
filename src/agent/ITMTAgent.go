package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/src/env"
)

type ITMTAgent interface {
	agent.IAgent[ITMTAgent]

	BroadcastPosition()

	TestMySynapse()

	PlayTurn()

	GetPos() env.Position

	ChangeEnergy(energyDelta int)

	GetEnergy() int

	ClearObstructions()

	NewObstructionMessage(pos env.Position) *ObstructionMessage
	HandleObstructionMessage(*ObstructionMessage)
}
