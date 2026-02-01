package env

type IEnvironment interface {
	GetTile(pos Position) (*Tile, bool)
	BoundPos(pos Position) Position
}
