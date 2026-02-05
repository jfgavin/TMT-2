package env

import (
	"math/rand"

	"github.com/google/uuid"
)

type Tile struct {
	Resources     int
	contributions map[uuid.UUID]int
}

func NewTile() *Tile {
	return &Tile{
		Resources:     0,
		contributions: make(map[uuid.UUID]int),
	}
}

func (tile *Tile) GetResources() int {
	sum := 0
	for _, res := range tile.contributions {
		sum += res
	}
	return sum
}

func (tile *Tile) GetContributions(id uuid.UUID) (int, bool) {
	cont, ok := tile.contributions[id]
	return cont, ok
}

func (tile *Tile) RefreshResources() {
	sum := 0
	for id, amt := range tile.contributions {
		if amt > 0 {
			sum += amt
		} else {
			// If no resources contributed by cluster id, delete the contribution entry
			delete(tile.contributions, id)
		}
	}
	tile.Resources = sum
}

func (tile *Tile) AddResources(source uuid.UUID, amt int) {
	tile.contributions[source] += amt
	tile.RefreshResources()
}

func (tile *Tile) SubResources(source uuid.UUID, amt int) {
	tile.contributions[source] -= amt
	tile.RefreshResources()
}

func (tile *Tile) DrainResources(amt int) {
	uuids := make([]uuid.UUID, 0)
	for uuid := range tile.contributions {
		uuids = append(uuids, uuid)
	}
	for amt > 0 && len(uuids) > 0 {
		index := rand.Intn(len(uuids))
		uuid := uuids[index]
		if tile.contributions[uuid] > 0 {
			tile.contributions[uuid]--
			amt--
		} else {
			// Cut this uuid out of selection
			uuids = append(uuids[:index], uuids[index+1:]...)
		}
	}
	tile.RefreshResources()
}
