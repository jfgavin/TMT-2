package agent

func (tmta *TMTAgent) HarvestResources() bool {
	tile, found := tmta.env.GetTile(tmta.Pos)
	if !found {
		return false
	}

	yield := tmta.cfg.ResourceYield
	available := tile.GetResources()

	if available <= 0 {
		return false
	} else if available < yield {
		// Add all available resources to agent as energy, and subtract from the tile
		tile.DrainResources(available)
		tmta.ChangeEnergy(available)
	} else {
		tile.DrainResources(yield)
		tmta.ChangeEnergy(yield)
	}

	return true
}
