package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/agent"
	"github.com/jfgavin/TMT-2/src/env"
)

type GameState struct {
	Iteration int
	Turn      int
	Grid      [][]*env.Tile
	Agents    map[uuid.UUID]agent.ITMTAgent
}

func BuildGameState(serv *GameServer, iteration, turn int) GameState {
	return GameState{
		Iteration: iteration,
		Turn:      turn,
		Grid:      serv.Env.Grid,
		Agents:    serv.GetAgentMap(),
	}
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
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to Python: %w", err)
	} else {
		fmt.Printf("Socket successfully initialised at %s\n", address)
	}
	serv.Conn = conn
	return nil
}

func (serv *GameServer) CloseSocket() {
	if serv.Conn != nil {
		serv.Conn.Close()
	}
}
