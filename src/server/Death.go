package server

import (
	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/agent"
)

type DeathType int

const (
	Elimination DeathType = iota
	Sacrifice
)

type DeathRecord struct {
	id        uuid.UUID
	parent    uuid.UUID
	deathType DeathType
}

func (serv *GameServer) AddDeathRecord(ag agent.ITMTAgent, dt DeathType) {
	newRecord := DeathRecord{
		ag.GetID(),
		ag.GetParent(),
		dt,
	}
	serv.deathRecords = append(serv.deathRecords, newRecord)
}
