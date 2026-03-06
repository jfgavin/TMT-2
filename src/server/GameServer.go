package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/google/uuid"
	"github.com/jfgavin/TMT-2/src/agent"
	"github.com/jfgavin/TMT-2/src/config"
	"github.com/jfgavin/TMT-2/src/env"
)

type GameServer struct {
	*server.BaseServer[agent.ITMTAgent]
	cfg          config.ServerConfig
	agCfg        config.AgentConfig
	Env          *env.Environment
	Conn         net.Conn
	obstructions map[env.Position]struct{}

	sacrificeReqs []uuid.UUID
	deathRecords  []DeathRecord
}

func (serv *GameServer) RunTurn(i, j int) {
	serv.UpdateObstructions()
	for _, ag := range serv.GetShuffledAgents() {
		ag.PlayTurn()
	}
	serv.StreamGameIteration()
	serv.HandleAgentMortality()
}

func (serv *GameServer) RunStartOfIteration(int) {
}

func (serv *GameServer) RunEndOfIteration(int) {
}

func (serv *GameServer) Start() {
	serv.BaseServer.Start()

	fmt.Println("Simulation Over!")
	serv.CloseSocket()
}

func NewGameServer(cfg config.Config) *GameServer {
	serv := &GameServer{
		BaseServer:   server.CreateBaseServer[agent.ITMTAgent](cfg.Serv.Iterations, cfg.Serv.Turns, 10*time.Millisecond, 100), // embed BaseServer: maxTimeout = 10ms, maxThreads = 100
		Env:          env.NewEnvironment(cfg.Env),
		cfg:          cfg.Serv,
		agCfg:        cfg.Agent, // Stored for spawning more agents later
		obstructions: make(map[env.Position]struct{}),
		deathRecords: make([]DeathRecord, 0),
	}

	// Add agents
	serv.IntroduceAgents()

	// set GameRunner to bind RunTurn to BaseServer
	serv.SetGameRunner(serv)

	// Initialise socket
	if err := serv.InitSocket("127.0.0.1:5000", cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Socket init failed: %v\n", err)
	}

	return serv
}
