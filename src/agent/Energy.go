package agent

func (tmta *TMTAgent) HarvestResources() bool {
	pos := tmta.Pos
	available, ok := tmta.env.GetResources()[pos]
	if !ok {
		return false
	}
	yield := tmta.cfg.ResourceYield

	if available < yield {
		// Add all available resources to agent as energy, and subtract from the tile
		tmta.env.DrainResources(pos, available)
		tmta.ChangeEnergy(available)
	} else {
		tmta.env.DrainResources(pos, yield)
		tmta.ChangeEnergy(yield)
	}

	return true
}
