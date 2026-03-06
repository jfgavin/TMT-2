package agent

import "github.com/google/uuid"

func (tmta *TMTAgent) AddChildID(id uuid.UUID) {
	tmta.Children = append(tmta.Children, id)
}

func (tmta *TMTAgent) GetChildren() []uuid.UUID {
	return tmta.Children
}

func (tmta *TMTAgent) SpawnChild() bool {
	ok := tmta.serv.RequestChild(tmta)
	return ok
}
