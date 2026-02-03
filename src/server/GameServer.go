package server

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/jfgavin/TMT-2/src/agent"
	"github.com/jfgavin/TMT-2/src/config"
	"github.com/jfgavin/TMT-2/src/env"
)

type GameServer struct {
	*server.BaseServer[agent.ITMTAgent]
	cfg  config.ServerConfig
	Env  *env.Environment
	Conn net.Conn
}

func (serv *GameServer) RunTurn(i, j int) {
	serv.ElimDrainedAgents()
	for _, ag := range serv.GetAgentMap() {
		ag.BroadcastPosition()
	}
	for _, ag := range serv.GetShuffledAgents() {
		ag.PlayTurn()
	}
	StreamGameIteration(serv, i, j)
	serv.DrainAgents()
	for _, ag := range serv.GetAgentMap() {
		ag.ClearObstructions()
	}
}

func (serv *GameServer) RunStartOfIteration(int) {
	serv.Env.IntroduceResources()
}

func (serv *GameServer) RunEndOfIteration(int) {

}

func (serv *GameServer) Start() {
	serv.BaseServer.Start()

	// Post-game functionality can go here
	fmt.Println("Game Over!")
	serv.CloseSocket()
}

func NewGameServer(cfg config.Config) *GameServer {
	serv := &GameServer{
		BaseServer: server.CreateBaseServer[agent.ITMTAgent](cfg.Serv.Iterations, cfg.Serv.Turns, 10*time.Millisecond, 100), // embed BaseServer: maxTimeout = 10ms, maxThreads = 100
		Env:        env.NewEnvironment(cfg.Env),
		cfg:        cfg.Serv,
	}

	for i := 0; i < cfg.Serv.NumAgents; i++ {
		pos := env.Position{X: rand.Intn(cfg.Env.GridSize), Y: rand.Intn(cfg.Env.GridSize)}

		ga := agent.NewTMTAgent(serv, cfg.Agent, serv.Env, fmt.Sprintf("Agent %d", i), pos)
		serv.AddAgent(ga)
	}

	// set GameRunner to bind RunTurn to BaseServer
	serv.SetGameRunner(serv)

	// Initialise socket
	if err := serv.InitSocket("127.0.0.1:5000"); err != nil {
		fmt.Fprintf(os.Stderr, "Socket init failed: %v\n", err)
	}

	return serv
}
