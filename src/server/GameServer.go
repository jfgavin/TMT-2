package server

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/jfgavin/TMT-2/src/config"
	"github.com/jfgavin/TMT-2/src/infra"
)

type GameServer struct {
	*server.BaseServer[infra.ITMTAgent]
	cfg  config.ServerConfig
	Env  *infra.Environment
	Conn net.Conn
}

func (serv *GameServer) RunTurn(i, j int) {
	for _, ag := range serv.GetAgentMap() {

		pos := ag.GetPos()

		pos.X = pos.X + (-1 + rand.Intn(3))
		pos.Y = pos.Y + (-1 + rand.Intn(3))

		ag.SetPos(pos)
	}
	StreamGameIteration(serv, i, j)
}

func (serv *GameServer) RunStartOfIteration(int) {
	for _, ag := range serv.GetAgentMap() {
		ag.ResetEnergy()
	}
	serv.Env.IntroduceResources(10000, 10)
}

func (serv *GameServer) RunEndOfIteration(int) {
	serv.RunMessagingTurn()
}

// make all agents message
func (serv *GameServer) RunMessagingTurn() {
	// Let agents signal they're ready for communication
	for _, gc := range serv.GetAgentMap() {
		gc.DoMessaging()
	}
}

// Socket
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

func (serv *GameServer) Start() {
	serv.BaseServer.Start()

	// Post-game functionality can go here
	fmt.Println("Game Over!")
	serv.CloseSocket()
}

func NewGameServer(cfg config.Config) *GameServer {
	serv := &GameServer{
		BaseServer: server.CreateBaseServer[infra.ITMTAgent](cfg.Serv.Iterations, cfg.Serv.Turns, 10*time.Millisecond, 100), // embed BaseServer: maxTimeout = 10ms, maxThreads = 100
		Env:        infra.NewEnvironment(cfg.Env),
		cfg:        cfg.Serv,
	}

	for i := 0; i < cfg.Serv.NumAgents; i++ {
		pos := infra.Position{X: rand.Intn(cfg.Env.GridSize), Y: rand.Intn(cfg.Env.GridSize)}

		ga := infra.NewTMTAgent(serv, cfg.Agent, serv.Env, fmt.Sprintf("Agent %d", i), pos)
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
