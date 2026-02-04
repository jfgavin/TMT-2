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

func (tmta *TMTAgent) IsReachable(target env.Position) bool {
	path := tmta.Pos.GreedyPath(target)

	// Walk path and check unobstructed
	for _, step := range path {
		if step.IsObstructed(tmta.obstructions) {
			return false
		}
	}
	return true
}

// Random move to one of the unobstructed adjascent cells, if possible
func (tmta *TMTAgent) GetRandomStep() (env.Position, bool) {
	adj := tmta.Pos.GetAdjacent()

	// Shuffle adjascent positions
	rand.Shuffle(len(adj), func(i, j int) {
		adj[i], adj[j] = adj[j], adj[i]
	})

	for _, pos := range adj {
		if !pos.IsObstructed(tmta.obstructions) && pos.IsBounded(tmta.env.GridSize()) {
			return pos, true
		}
	}

	return tmta.Pos, false
}

// Returns reachable position with highest utility (resources / dist + 1)
func (tmta *TMTAgent) GetBestTarget() (env.Position, bool) {
	locals := tmta.VisiblePositions()
	startPos := tmta.Pos

	bestTarget := startPos
	bestUtility := 0.0

	for _, target := range locals {
		tile, ok := tmta.env.GetTile(target)
		if !ok {
			continue
		}
		dist := startPos.ManhatDist(target)
		tileUtility := float64(tile.GetResources()) / float64(dist+1)
		if tileUtility > bestUtility && tmta.IsReachable(target) {
			bestTarget = target
			bestUtility = tileUtility
		}
	}

	return bestTarget, bestUtility > 0
}

// Try to move to resources, otherwise explore, otherwise stand still
func (tmta *TMTAgent) Move() {
	tmta.Target = env.Position{}
	step := tmta.Pos

	if best, ok := tmta.GetBestTarget(); ok {
		step = tmta.Pos.GreedyNextStep(best)
		tmta.Target = best
	} else if randStep, ok := tmta.GetRandomStep(); ok {
		step = randStep
	} else {
		return
	}

	tmta.SetPosAndBroadcast(step)
}
