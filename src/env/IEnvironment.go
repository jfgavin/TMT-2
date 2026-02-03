package env

type IEnvironment interface {
	GridSize() int
	GetGrid() [][]*Tile
	GetTile(pos Position) (*Tile, bool)
	TilePos(tile *Tile) Position
	LocalTiles(tile *Tile, radius int) []*Tile
	BoundPos(pos Position) Position
}
