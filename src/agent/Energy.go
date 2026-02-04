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
		tmta.env.ChangeResources(tmta.Pos, -available)
		tmta.ChangeEnergy(available)
	} else {
		tmta.env.ChangeResources(tmta.Pos, -yield)
		tmta.ChangeEnergy(yield)
	}

	return true
}
