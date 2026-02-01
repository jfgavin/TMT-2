package agent

import (
	"math/rand"

	"github.com/jfgavin/TMT-2/src/env"
)

func (tmta *TMTAgent) getPercievedTiles() []*env.Tile {
	percievedTiles := make([]*env.Tile, 0)

	for y, row := range tmta.env.GetGrid() {
		for x, tile := range row {
			tilePos := env.Position{X: x, Y: y}
			if tmta.Pos.ManhatDist(tilePos) <= tmta.cfg.PerceptiveRange {
				percievedTiles = append(percievedTiles, tile)
			}
		}
	}

	return percievedTiles
}

func (tmta *TMTAgent) getRandAdjTile() *env.Tile {
	adjPositions := tmta.Pos.GetAdjacent()

	validTiles := make([]*env.Tile, 0, len(adjPositions))
	for _, pos := range adjPositions {
		if tile, found := tmta.env.GetTile(pos); found {
			validTiles = append(validTiles, tile)
		}
	}

	if len(validTiles) == 0 {
		return nil
	}

	return validTiles[rand.Intn(len(validTiles))]
}

func (tmta *TMTAgent) getTargetTile() *env.Tile {
	bestTile := tmta.getRandAdjTile()
	bestUtility := 0
	for _, tile := range tmta.getPercievedTiles() {
		utility := tile.Resources
		if utility > bestUtility {
			bestUtility = utility
			bestTile = tile
		}
	}

	return bestTile
}

func (tmta *TMTAgent) getTargetPos() env.Position {
	return tmta.env.TilePos(tmta.getTargetTile())
}

func (tmta *TMTAgent) Move() {
	targetPos := tmta.getTargetPos()
	nextPos := tmta.Pos.GetNextStep(targetPos)
	tmta.Pos = nextPos
}
