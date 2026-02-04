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

	return path
}

// Greedily get next position which reduces Manhattan distance to target
func (pos Position) GreedyNextStep(target Position) Position {
	nextStep := pos
	remDist := pos.ManhatDist(target)
	bestRemDist := remDist

	for _, adj := range pos.GetAdjacent() {
		dist := adj.ManhatDist(target)
		if dist < bestRemDist {
			nextStep = adj
			bestRemDist = dist
		}
	}
	return nextStep
}

// Concatenate greedy steps to make a path to target
func (pos Position) GreedyPath(target Position) []Position {
	dist := pos.ManhatDist(target)
	path := make([]Position, 0, dist)

	current := pos
	path = append(path, current)

	for current != target {
		current = current.GreedyNextStep(target)
		path = append(path, current)
	}

	return path
}

func (pos Position) IsObstructed(obstructions map[Position]struct{}) bool {
	_, blocked := obstructions[pos]
	return blocked
}
