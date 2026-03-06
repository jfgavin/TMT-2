package agent

import (
	"github.com/jfgavin/TMT-2/src/env"
)

// Very simple pathfinding algo, where valid step to reduce distance to target is chosen
func (tmta *TMTAgent) GetGreedyPath(target env.Position) ([]env.Position, bool) {
	current := tmta.Pos
	if current == target {
		return nil, false
	}

	remainingDist := current.ManhatDist(target)

	path := make([]env.Position, 0, remainingDist)
	for current != target {
		foundNextStep := false
		for _, adj := range current.GetShuffledAdjacent() {
			dist := adj.ManhatDist(target)
			if dist < remainingDist && !tmta.serv.IsObstructed(adj) {
				current = adj
				remainingDist = dist
				foundNextStep = true
				break
			}
		}
		if !foundNextStep {
			return path, false
		}
		path = append(path, current)
	}

	return path, true
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

// Random move to one of the unobstructed adjascent cells, if possible
func (tmta *TMTAgent) GetRandomStep() (env.Position, bool) {
	for _, adj := range tmta.Pos.GetShuffledAdjacent() {
		if !tmta.serv.IsObstructed(adj) {
			return adj, true
		}
	}
	return tmta.Pos, false
}

// Returns reachable position with highest utility (resources / dist + 1)
func (tmta *TMTAgent) GetBestStep() (env.Position, bool) {
	startPos := tmta.Pos
	bestStep := startPos
	bestUtility := 0.0

	resourceMap := tmta.env.GetResources()
	for _, pos := range tmta.VisiblePositions() {
		resources, ok := resourceMap[pos]
		if !ok || resources <= 0 {
			continue
		}

		tileUtility := float64(resources) / float64(startPos.ManhatDist(pos)+1)
		if tileUtility <= bestUtility {
			continue
		}

		if pos == startPos {
			bestStep = startPos
			bestUtility = tileUtility
			continue
		}

		path, ok := tmta.GetGreedyPath(pos)
		if ok {
			bestStep = path[0]
			bestUtility = tileUtility
		}
	}
	return bestStep, bestUtility > 0.0
}

func (tmta *TMTAgent) Move() bool {
	step := tmta.Pos

	if bestStep, ok := tmta.GetBestStep(); ok {
		step = bestStep
	} else if randStep, ok := tmta.GetRandomStep(); ok {
		step = randStep
	}

	return tmta.serv.MoveAgent(tmta, step)
}
