package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/jfgavin/TMT-2/src/config"
	"github.com/jfgavin/TMT-2/src/infra"
)

type GameServer struct {
	*server.BaseServer[infra.IGameAgent]
	Env  *infra.Environment
	Conn net.Conn
	cfg  config.Config
}

func (serv *GameServer) RunTurn(i, j int) {
	StreamGameIteration(serv, i, j)
}

func (serv *GameServer) RunStartOfIteration(int) {
	for _, ag := range serv.GetAgentMap() {
		ag.SetEnergy(serv.cfg.StartingEnergy)
	}
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

func NewGameServer(serverConfig config.Config) *GameServer {
	serv := &GameServer{
		BaseServer: server.CreateBaseServer[infra.IGameAgent](serverConfig.Iterations, serverConfig.Turns, 10*time.Millisecond, 100), // embed BaseServer: maxTimeout = 10ms, maxThreads = 100
		Env:        infra.NewEnvironment(serverConfig.GridSize),
		cfg:        serverConfig,
	}

	// set GameRunner to bind RunTurn to BaseServer
	serv.SetGameRunner(serv)

	// Initialise socket
	if err := serv.InitSocket("127.0.0.1:5000"); err != nil {
		fmt.Fprintf(os.Stderr, "Socket init failed: %v\n", err)
	}

	return serv
}
