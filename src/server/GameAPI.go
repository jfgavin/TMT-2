package server

import (
	"encoding/json"
	"fmt"

	"github.com/jfgavin/TMT-2/src/infra"
)

type GameState struct {
	Iteration int
	Turn      int
	Grid      [][]*infra.Tile
}

func BuildGameState(serv *GameServer, iteration, turn int) GameState {
	return GameState{
		Iteration: iteration,
		Turn:      turn,
		Grid:      serv.Env.Grid,
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
