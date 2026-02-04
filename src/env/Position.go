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

func (pos Position) IsBounded(upperBound int) bool {
	return pos.X >= 0 && pos.X < upperBound && pos.Y >= 0 && pos.Y < upperBound
}

func (pos Position) Bound(upperBound int) {
	if pos.X < 0 {
		pos.X = 0
	} else if pos.X > upperBound {
		pos.X = upperBound
	}

	if pos.Y < 0 {
		pos.Y = 0
	} else if pos.Y > upperBound {
		pos.Y = upperBound
	}
}

func (pos Position) GetAdjacent() [4]Position {
	return [4]Position{
		{pos.X + 1, pos.Y},
		{pos.X - 1, pos.Y},
		{pos.X, pos.Y + 1},
		{pos.X, pos.Y - 1},
	}
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

func (pos Position) IsObstructed(obstructions map[Position]struct{}) bool {
	_, blocked := obstructions[pos]
	return blocked
}

func (pos Position) UnobstructedNextStep(target Position, obstructions map[Position]struct{}) Position {
	for _, step := range pos.GetScoredNextSteps(target) {
		if !step.IsObstructed(obstructions) {
			return step
		}
	}
	return pos // No step
}
