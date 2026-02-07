package env

type IEnvironment interface {
	GridSize() int
	GetTile(pos Position) (*Tile, bool)
}
