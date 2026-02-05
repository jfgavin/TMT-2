package agent

import (
	"github.com/jfgavin/TMT-2/src/env"
)

// Update position, and broadcast this new obstruction to all other agents
func (tmta *TMTAgent) SetPosAndBroadcast(pos env.Position) {
	tmta.Pos = pos
	tmta.BroadcastPosition()
}

// Returns all positions with Manhattan distance <= agent's visual range
func (tmta *TMTAgent) VisiblePositions() []env.Position {
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

// Checks if target is reachable, then returns next step if true
func (tmta *TMTAgent) IsReachable(target env.Position) (env.Position, bool) {
	// If getting there and harvesting is more energy than agent has, then false
	if dist := tmta.Pos.ManhatDist(target); dist+1 > tmta.Energy {
		return env.Position{}, false
	}

	path := tmta.Pos.GreedyPath(target)

	// Walk path and check unobstructed
	for _, step := range path {
		if step.IsObstructed(tmta.obstructions) {
			return env.Position{}, false
		}
	}

	return path[min(1, len(path))], true
}

// Random move to one of the unobstructed adjascent cells, if possible
func (tmta *TMTAgent) GetRandomStep() (env.Position, bool) {
	adj := tmta.Pos.GetAdjacent()

	for _, pos := range adj {
		if !pos.IsObstructed(tmta.obstructions) && pos.IsBounded(tmta.env.GridSize()) {
			return pos, true
		}
	}

	return tmta.Pos, false
}

// Returns reachable position with highest utility (resources / dist + 1)
func (tmta *TMTAgent) GetBestStep() (env.Position, bool) {
	locals := tmta.VisiblePositions()
	startPos := tmta.Pos

	bestStep := startPos
	bestUtility := 0.0

	for _, target := range locals {
		tile, ok := tmta.env.GetTile(target)
		if !ok || tile.GetResources() < 0 {
			continue
		}
		dist := startPos.ManhatDist(target)
		tileUtility := float64(tile.GetResources()) / float64(dist+1)
		if tileUtility > bestUtility {
			if step, ok := tmta.IsReachable(target); ok {
				bestStep = step
				bestUtility = tileUtility
			}
		}

	}
	return bestStep, bestUtility > 0
}

// Try to move to resources, otherwise explore, otherwise stand still
func (tmta *TMTAgent) Move() {
	tmta.Target = env.Position{}
	step := tmta.Pos

	if bestStep, ok := tmta.GetBestStep(); ok {
		step = bestStep
	} else if randStep, ok := tmta.GetRandomStep(); ok {
		step = randStep
	} else {
		return
	}

	tmta.SetPosAndBroadcast(step)
}
