package agent

import (
	"math/rand/v2"

	"github.com/jfgavin/TMT-2/src/env"
)

// Update position, and broadcast this new obstruction to all other agents
func (tmta *TMTAgent) SetPosAndBroadcast(pos env.Position) {
	tmta.Pos = pos
	tmta.BroadcastPosition()
}

// Returns all positions with Manhattan distance <= agent's visual range
func (tmta *TMTAgent) GetVisible() []env.Position {
	pos := tmta.Pos
	visMax := tmta.cfg.VisualRange

	capacity := 1 + 2*visMax*(visMax+1)
	out := make([]env.Position, 0, capacity)

	for dy := -visMax; dy <= visMax; dy++ {
		limit := visMax - max(dy, -dy)
		for dx := -limit; dx <= limit; dx++ {
			local := env.Position{
				X: pos.X + dx,
				Y: pos.Y + dy,
			}
			if local.IsBounded(tmta.env.GridSize()) {
				out = append(out, local)
			}
		}
	}
	return out
}

// Random move to one of the unobstructed adjascent cells, if possible
func (tmta *TMTAgent) ExploreMove() {
	adj := tmta.Pos.GetAdjacent()

	// Shuffle adjascent positions
	rand.Shuffle(len(adj), func(i, j int) {
		adj[i], adj[j] = adj[j], adj[i]
	})

	for _, pos := range adj {
		if !pos.IsObstructed(tmta.obstructions) && pos.IsBounded(tmta.env.GridSize()) {
			tmta.SetPosAndBroadcast(pos)
			return
		}
	}
}

// Move towards the unobstructed cell with most resources
func (tmta *TMTAgent) TargetedMove() {
	locals := tmta.GetVisible()

	startPos := tmta.Pos
	bestTarget := startPos
	bestUtility := 0

	for _, target := range locals {
		tile, ok := tmta.env.GetTile(target)
		if !ok {
			return
		}
		dist := startPos.ManhatDist(target)
		tileUtility := tile.GetResources() / (dist + 1)
		if tileUtility > bestUtility && !target.IsObstructed(tmta.obstructions) {
			bestTarget = target
			bestUtility = tileUtility
		}
	}

	// If the best unobstructed position has resources, move a step towards it, else, explore
	if bestUtility > 0 {
		nextStep := startPos.UnobstructedNextStep(bestTarget, tmta.obstructions)
		tmta.SetPosAndBroadcast(nextStep)
	} else {
		tmta.ExploreMove()
	}
}
