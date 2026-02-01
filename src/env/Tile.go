package env

type Tile struct {
	Resources int
}

func NewTile(resources int) *Tile {
	return &Tile{
		Resources: resources,
	}
}

func (tile *Tile) GetResources() int {
	return tile.Resources
}

func (tile *Tile) ChangeResources(resourceDelta int) {
	tile.Resources += resourceDelta
}
