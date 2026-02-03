package env

import "sort"

type Position struct {
	X, Y int
}

func (posA *Position) ManhatDist(posB Position) int {
	x := posA.X - posB.X
	y := posA.Y - posB.Y

	x_abs := max(x, -x)
	y_abs := max(y, -y)

	return x_abs + y_abs
}

func (pos *Position) LocalPositions(maxDist int) []Position {
	// Returns circle of positions within maxDist Manhattan distance
	localPositions := make([]Position, 0)
	for y := -maxDist; y <= maxDist; y++ {
		for x := -maxDist; x <= maxDist; x++ {
			localPos := Position{X: x, Y: y}
			if pos.ManhatDist(localPos) <= maxDist {
				localPositions = append(localPositions, localPos)
			}
		}
	}
	return localPositions
}

func (pos *Position) GetAdjacent() []Position {
	// Returns list of adjascent positions
	return []Position{
		{
			X: pos.X + 1,
			Y: pos.Y,
		},
		{
			X: pos.X - 1,
			Y: pos.Y,
		},
		{
			X: pos.X,
			Y: pos.Y + 1,
		},
		{
			X: pos.X,
			Y: pos.Y - 1,
		},
	}
}

func (pos *Position) GetScoredNextSteps(target Position) []Position {
	// Store baseline (waiting)
	baseScore := pos.ManhatDist(target)

	type scoredStep struct {
		pos   Position
		score int
	}

	steps := make([]scoredStep, 0, 5)

	// Waiting is always a candidate
	steps = append(steps, scoredStep{
		pos:   *pos,
		score: baseScore,
	})

	// Adjacent moves
	for _, adj := range pos.GetAdjacent() {
		score := adj.ManhatDist(target)

		// Drop anything worse than waiting
		if score <= baseScore {
			steps = append(steps, scoredStep{
				pos:   adj,
				score: score,
			})
		}
	}

	// Sort best → worst
	sort.Slice(steps, func(i, j int) bool {
		if steps[i].score != steps[j].score {
			return steps[i].score < steps[j].score
		}
		// Stable tie-breaker (important to avoid oscillation)
		if steps[i].pos.X != steps[j].pos.X {
			return steps[i].pos.X < steps[j].pos.X
		}
		return steps[i].pos.Y < steps[j].pos.Y
	})

	// Strip scores
	result := make([]Position, len(steps))
	for i, s := range steps {
		result[i] = s.pos
	}

	return result
}
