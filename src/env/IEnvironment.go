package env

type IEnvironment interface {
	GridSize() int
	GetResources() map[Position]int
	DrainResources(pos Position, amt int) bool
	PlaceTombstone(pos Position)
	PlaceMemorial(pos Position)
	TickGraves()
	GetGraves() map[Position]*Grave
}
