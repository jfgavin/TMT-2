package agent

import (
	"github.com/jfgavin/TMT-2/src/env"
)

func (tmta *TMTAgent) GetTargetPos() env.Position {
	visRange := tmta.cfg.VisualRange
	locals := tmta.Pos.LocalPositions(visRange)

	bestPos := tmta.Pos
	bestResources := 0

	for _, pos := range locals {
		tile, ok := tmta.env.GetTile(pos)
		if ok && tile.Resources > bestResources {
			bestPos = pos
		}
	}

	return bestPos
}

func (tmta *TMTAgent) Move() {
	tmta.Pos = tmta.Pos.UnobstructedNextStep(tmta.GetTargetPos(), tmta.obstructions)
	tmta.BroadcastPosition()
}

func (tmta *TMTAgent) HandleObstructionMessage(msg *ObstructionMessage) {
	tmta.obstructions[msg.Pos] = struct{}{}
}
