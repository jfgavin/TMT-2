package env

type Position struct {
	X, Y int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (a Position) ManhatDist(b Position) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

// Diamond (Manhattan ball) around pos
func (pos Position) LocalPositions(maxDist int) []Position {
	capacity := 1 + 2*maxDist*(maxDist+1)
	out := make([]Position, 0, capacity)

	for dy := -maxDist; dy <= maxDist; dy++ {
		limit := maxDist - abs(dy)
		for dx := -limit; dx <= limit; dx++ {
			local := Position{
				X: pos.X + dx,
				Y: pos.Y + dy,
			}
			out = append(out, local)
		}
	}
	return out
}

func (pos Position) GetAdjacent() [4]Position {
	return [4]Position{
		{pos.X + 1, pos.Y},
		{pos.X - 1, pos.Y},
		{pos.X, pos.Y + 1},
		{pos.X, pos.Y - 1},
	}

	return path
}

func (pos Position) GetScoredNextSteps(target Position) []Position {
	base := pos.ManhatDist(target)

	// Worst case: No step + 4 adj
	out := make([]Position, 0, 5)

	bestScore := base

	// Always allow no step
	out = append(out, pos)

	for _, adj := range pos.GetAdjacent() {
		score := adj.ManhatDist(target)

		if score < bestScore {
			// Found strictly better move, so reset list
			bestScore = score
			out = out[:0]
			out = append(out, adj)
		} else if score == bestScore {
			out = append(out, adj)
		}
	}

	return out
}

func (pos Position) UnobstructedNextStep(target Position, obstructions map[Position]struct{}) Position {
	for _, step := range pos.GetScoredNextSteps(target) {
		if _, blocked := obstructions[step]; !blocked {
			return step
		}
	}
	return pos // No step
}
