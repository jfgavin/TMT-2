package env

type IEnvironment interface {
	GridSize() int
	GetTile(pos Position) (Tile, bool)
	ChangeResources(pos Position, delta int) bool
	RandomlyAddResources(amt int)
}
