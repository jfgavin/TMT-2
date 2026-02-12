package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/agent"
	"github.com/jfgavin/TMT-2/src/env"
)

type ResourceEntry struct {
	Pos   env.Position
	Value int
}

type GraveEntry struct {
	Pos   env.Position
	Grave *env.Grave
}

type Metadata struct {
	GridSize int
}
type GameState struct {
	Iteration int
	Turn      int
	Resources []ResourceEntry
	Graves    []GraveEntry
	Agents    map[uuid.UUID]agent.ITMTAgent
}

func BuildGameState(serv *GameServer, iteration, turn int) GameState {
	gs := GameState{
		Iteration: iteration,
		Turn:      turn,
		Resources: make([]ResourceEntry, 0),
		Graves:    make([]GraveEntry, 0),
		Agents:    serv.GetAgentMap(),
	}

	for pos, val := range serv.Env.GetResources() {
		gs.Resources = append(gs.Resources, ResourceEntry{Pos: pos, Value: val})
	}

	for pos, grv := range serv.Env.GetGraves() {
		gs.Graves = append(gs.Graves, GraveEntry{Pos: pos, Grave: grv})
	}

	return gs
}

func StreamGameIteration(serv *GameServer, iteration, turn int) error {
	if serv.Conn == nil {
		return fmt.Errorf("no connection")
	}

	state := BuildGameState(serv, iteration, turn)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	_, err = serv.Conn.Write(append(data, '\n'))
	return err
}

// Websocket
func (serv *GameServer) InitSocket(address string) error {
	fmt.Println("Connecting...")
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to Python: %w", err)
	} else {
		fmt.Printf("Socket successfully initialised at %s\n", address)
	}
	serv.Conn = conn

	metadata, err := json.Marshal(Metadata{GridSize: serv.Env.GridSize()})
	if err != nil {
		return err
	}
	_, err = serv.Conn.Write(append(metadata, '\n'))
	return err
}

func (serv *GameServer) CloseSocket() {
	if serv.Conn != nil {
		serv.Conn.Close()
	}
}
