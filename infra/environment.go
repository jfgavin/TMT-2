package infra

import (
	"math/rand"
	"time"
)

const (
	NumAgents          = 4
	StartingEnergy     = 5
	StartingGridWidth  = 8
	StartingGridHeight = 8
)

// CreateEnvironment initializes the game environment with mall tiles, exits, and loot shops
// heroNames: list of hero names to assign shops to (each hero gets one shop)
func CreateEnvironment(width, height int, heroNames []string) *Environment {
	env := &Environment{
		Grid:        make(map[Position]Tile),
		Exits:       []Position{},
		Shops:       []Position{},
		LootedShops: make(map[Position]bool),
	}

	// Initialize all tiles as normal
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := Position{X: x, Y: y}
			env.Grid[pos] = Tile{Type: NormalTile}
		}
	}

	// Place exits at corners (0,0), (width-1,0), (0,height-1), (width-1,height-1)
	exits := []Position{
		{X: 0, Y: 0},
		{X: width - 1, Y: 0},
		{X: 0, Y: height - 1},
		{X: width - 1, Y: height - 1},
	}
	for _, exit := range exits {
		env.Grid[exit] = Tile{Type: ExitTile}
		env.Exits = append(env.Exits, exit)
	}

	// Place loot shops - one for each hero, at random positions
	// Shops cannot be placed on exits, and must be unique
	usedPositions := make(map[Position]bool)
	// Mark exit positions as used
	for _, exit := range exits {
		usedPositions[exit] = true
	}

	// Initialize random seed (use current time)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, heroName := range heroNames {
		// Try to find a valid random position for this shop
		var shop Position
		attempts := 0
		maxAttempts := width * height * 2 // Prevent infinite loop

		for attempts < maxAttempts {
			shop = Position{
				X: rng.Intn(width),
				Y: rng.Intn(height),
			}

			// Check if position is valid (not an exit and not already used)
			if !usedPositions[shop] {
				usedPositions[shop] = true
				// Create shop with owner
				env.Grid[shop] = Tile{Type: LootShopTile, Owner: heroName}
				env.Shops = append(env.Shops, shop)
				break
			}
			attempts++
		}

		// Fallback: if we couldn't find a random position, use a systematic approach
		if attempts >= maxAttempts {
			// Find first available non-exit position
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					shop = Position{X: x, Y: y}
					if !usedPositions[shop] {
						usedPositions[shop] = true
						env.Grid[shop] = Tile{Type: LootShopTile, Owner: heroName}
						env.Shops = append(env.Shops, shop)
						goto nextHero
					}
				}
			}
		}
	nextHero:
	}

	return env
}

// GetTile returns the tile at the given position
func (env *Environment) GetTile(pos Position) Tile {
	if tile, ok := env.Grid[pos]; ok {
		return tile
	}
	return Tile{Type: NormalTile}
}

// IsExit checks if the position is an exit
func (env *Environment) IsExit(pos Position) bool {
	return env.GetTile(pos).Type == ExitTile
}

// IsLootShop checks if the position is a loot shop
func (env *Environment) IsLootShop(pos Position) bool {
	return env.GetTile(pos).Type == LootShopTile
}

// IsMyLootShop checks if the position is a loot shop owned by the specified hero
func (env *Environment) IsMyLootShop(pos Position, heroName string) bool {
	tile := env.GetTile(pos)
	return tile.Type == LootShopTile && tile.Owner == heroName
}

// GetShopOwner returns the owner of a shop at the given position, or empty string
func (env *Environment) GetShopOwner(pos Position) string {
	return env.GetTile(pos).Owner
}

// MarkShopAsLooted marks a shop as looted
func (env *Environment) MarkShopAsLooted(pos Position) {
	env.LootedShops[pos] = true
}

// IsShopLooted checks if a shop has been looted
func (env *Environment) IsShopLooted(pos Position) bool {
	return env.LootedShops[pos]
}

// AllShopsLooted checks if all shops have been looted
func (env *Environment) AllShopsLooted() bool {
	for _, shopPos := range env.Shops {
		if !env.IsShopLooted(shopPos) {
			return false
		}
	}
	return true
}
