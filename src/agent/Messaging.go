package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/jfgavin/TMT-2/src/env"
)

type ObstructionMessage struct {
	message.BaseMessage
	Pos env.Position
}

func (tmta *TMTAgent) NewObstructionMessage(pos env.Position) *ObstructionMessage {
	return &ObstructionMessage{
		BaseMessage: tmta.CreateBaseMessage(),
		Pos:         pos,
	}
}

func (msg *ObstructionMessage) InvokeMessageHandler(agent ITMTAgent) {
	agent.HandleObstructionMessage(msg)
}
