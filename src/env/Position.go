package env

import "math/rand"

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

func (pos Position) GetShuffledAdjacent() [4]Position {
	adj := [4]Position{
		{pos.X + 1, pos.Y},
		{pos.X - 1, pos.Y},
		{pos.X, pos.Y + 1},
		{pos.X, pos.Y - 1},
	}

	rand.Shuffle(len(adj), func(i, j int) {
		adj[i], adj[j] = adj[j], adj[i]
	})

	return adj
}
