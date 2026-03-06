package server

import "github.com/jfgavin/TMT-2/src/agent"

func (serv *GameServer) GetElimCount() int {
	count := 0
	for _, rec := range serv.deathRecords {
		if rec.deathType == Elimination {
			count++
		}
	}
	return count
}

func (serv *GameServer) GetSacrificeCount() int {
	count := 0
	for _, rec := range serv.deathRecords {
		if rec.deathType == Sacrifice {
			count++
		}
	}
	return count
}

func (serv *GameServer) GetDeceasedChildCount(parent agent.ITMTAgent) int {
	parentID := parent.GetID()
	count := 0
	for _, rec := range serv.deathRecords {
		if rec.parent == parentID {
			count++
		}
	}
	return count
}
