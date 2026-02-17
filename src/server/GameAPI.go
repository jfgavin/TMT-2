package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/agent"
)

type Metadata struct {
	GridSize int
}
type GameState struct {
	Iteration int
	Turn      int
	Resources [][3]int
	Graves    [][3]any
	Agents    map[uuid.UUID]agent.ITMTAgent
}

func BuildGameState(serv *GameServer, iteration, turn int) GameState {
	resources := serv.Env.GetResources()
	graves := serv.Env.GetGraves()

	gs := GameState{
		Iteration: iteration,
		Turn:      turn,
		Resources: make([][3]int, 0, len(resources)),
		Graves:    make([][3]any, 0, len(graves)),
		Agents:    serv.GetAgentMap(),
	}

	for pos, val := range resources {
		gs.Resources = append(gs.Resources, [3]int{pos.X, pos.Y, val})
	}

	for pos, grv := range graves {
		gs.Graves = append(gs.Graves, [3]any{pos.X, pos.Y, grv})
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
