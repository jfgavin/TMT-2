package infra

type Position struct {
	X, Y int
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
	NoDirection
)

var Movement = map[Direction]Position{
	North:       {X: 0, Y: -1},
	South:       {X: 0, Y: 1},
	East:        {X: 1, Y: 0},
	West:        {X: -1, Y: 0},
	NoDirection: {X: 0, Y: 0},
}

var MovementName = map[Direction]string{
	North:       "North",
	South:       "South",
	East:        "East\t",
	West:        "West\t",
	NoDirection: "No Direction",
}
