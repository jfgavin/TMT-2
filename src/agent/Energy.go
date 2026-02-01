package agent

func (tmta *TMTAgent) HarvestResources() {
	tile, found := tmta.env.GetTile(tmta.Pos)
	if !found {
		return
	}

	yield := tmta.cfg.ResourceYield
	available := tile.GetResources()

	if available <= yield {
		// Add all available resources to agent as energy, and subtract from the tile
		tmta.ChangeEnergy(available)
		tile.ChangeResources(-available)
	} else {
		tmta.ChangeEnergy(yield)
		tile.ChangeResources(-yield)
	}
}
