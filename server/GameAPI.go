package gameServer

import (
	"encoding/json"
	"fmt"

	"github.com/jfgavin/TMT-2/infra"
)

// JSONTile represents a single tile for output
type JSONTile struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

// HeroStats for output
type HeroStats struct {
	Name   string `json:"name"`
	Energy int    `json:"energy"`
	Loot   bool   `json:"looted"`
}

// GameIteration represents a single frame/output
type GameIteration struct {
	Iteration int          `json:"iteration"`
	Turn      int          `json:"turn"`
	Board     [][]JSONTile `json:"board"`
	Heroes    []HeroStats  `json:"heroes"`
}

// StreamGameIteration outputs the current game state as JSON
func StreamGameIteration(serv *GameServer, iteration, turn int) error {
	// Check socket
	if serv.conn == nil {
		return fmt.Errorf("no active socket connection")
	}

	// Build tile grid
	gridLookup := make(map[infra.Position]infra.IGameAgent)
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			gridLookup[gc.GetPos()] = gc
		}
	}

	board := make([][]JSONTile, serv.GridHeight())
	for y := 0; y < serv.GridHeight(); y++ {
		row := make([]JSONTile, serv.GridWidth())
		for x := 0; x < serv.GridWidth(); x++ {
			pos := infra.Position{X: x, Y: y}
			if gc, ok := gridLookup[pos]; ok {
				heroChar := []rune(gc.GetName())[0]
				row[x] = JSONTile{
					Type: "hero",
					Name: string(heroChar),
				}
			} else {
				tile := serv.Environment().GetTile(pos)
				switch tile.Type {
				case infra.ExitTile:
					row[x] = JSONTile{Type: "exit"}
				case infra.LootShopTile:
					if serv.Environment().IsShopLooted(pos) {
						row[x] = JSONTile{Type: "looted_shop"}
					} else {
						row[x] = JSONTile{Type: "shop"}
					}
				default:
					row[x] = JSONTile{Type: "empty"}
				}
			}
		}
		board[y] = row
	}

	// Build hero stats
	var heroes []HeroStats
	for _, gc := range serv.GetAgentMap() {
		heroes = append(heroes, HeroStats{
			Name:   gc.GetName(),
			Energy: gc.GetEnergy(),
			Loot:   gc.HasLoot(),
		})
	}

	GameIteration := GameIteration{
		Iteration: iteration,
		Turn:      turn,
		Board:     board,
		Heroes:    heroes,
	}

	output, err := json.Marshal(GameIteration)
	if err != nil {
		return err
	}

	output = append(output, '\n')
	_, err = serv.conn.Write(output)
	return err
}
