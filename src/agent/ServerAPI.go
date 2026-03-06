package agent

import (
	"github.com/jfgavin/TMT-2/src/env"
)

type ServerAPI interface {
	GetElimCount() int
	RequestSacrifice(ITMTAgent)
	IsObstructed(env.Position) bool
	MoveAgent(ITMTAgent, env.Position) bool
	RequestChild(ITMTAgent) bool
}
