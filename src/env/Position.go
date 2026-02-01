package env

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

func (pos *Position) GetNextStep(target Position) Position {
	// Optimised next-step movement from pos to approach target
	step := *pos

	dx := target.X - pos.X
	dy := target.Y - pos.Y

	dxAbs := max(dx, -dx)
	dyAbs := max(dy, -dy)

	if dx == 0 && dy == 0 {
		return step
	}

	if dxAbs >= dyAbs {
		if dx > 0 {
			step.X++
		} else {
			step.X--
		}
	} else {
		if dy > 0 {
			step.Y++
		} else {
			step.Y--
		}
	}

	return step
}
