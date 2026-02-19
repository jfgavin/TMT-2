package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/agent"
)

type Message struct {
	Type string
	Data any
}

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

func writeMessage(conn net.Conn, tp string, v any) error {
	msg := Message{
		Type: tp,
		Data: v,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Write length prefix
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))

	_, err = conn.Write(lenBuf[:])
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func StreamGameIteration(serv *GameServer, iteration, turn int) error {
	if serv.Conn == nil {
		return fmt.Errorf("no connection")
	}

	state := BuildGameState(serv, iteration, turn)
	err := writeMessage(serv.Conn, "state", state)
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

	metadata := Metadata{
		GridSize: serv.Env.GridSize(),
	}

	err = writeMessage(serv.Conn, "metadata", metadata)
	return err
}

func (serv *GameServer) CloseSocket() {
	if serv.Conn != nil {
		serv.Conn.Close()
	}
}
