package agent

import "github.com/google/uuid"

func (tmta *TMTAgent) GetParent() uuid.UUID {
	return tmta.Parent
}

func (tmta *TMTAgent) SpawnChild() bool {
	ok := tmta.serv.RequestChild(tmta)
	return ok
}
