package infra

// TileType represents different types of tiles in the mall
type TileType int

const (
	NormalTile TileType = iota
	ExitTile
	LootShopTile
)

var TileTypeName = map[TileType]string{
	NormalTile:   "Normal",
	ExitTile:     "Exit",
	LootShopTile: "LootShop",
}

// Tile represents a tile in the mall grid
type Tile struct {
	Type  TileType
	Owner string // For LootShopTile: which hero owns this shop (empty for other tiles)
}

// Environment represents the game environment
type Environment struct {
	Grid        map[Position]Tile
	Exits       []Position
	Shops       []Position
	LootedShops map[Position]bool // Track which shops have been looted
}
