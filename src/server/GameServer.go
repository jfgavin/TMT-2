package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/jfgavin/TMT-2/src/infra"
)

type GameServer struct {
	*server.BaseServer[infra.IGameAgent]
	Env  *infra.Environment
	Conn net.Conn
}

func (serv *GameServer) RunTurn(i, j int) {
	StreamGameIteration(serv, i, j)
}

func (serv *GameServer) RunStartOfIteration(int) {
	for _, ag := range serv.GetAgentMap() {
		ag.SetEnergy(infra.StartingEnergy)
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

// override start
func (serv *GameServer) Start() {
	// steal method from package...
	serv.BaseServer.Start()

	// ...and add some more functionality for after the game
	fmt.Println("Game Over!")
	serv.CloseSocket()
}

// constructor for GameServer
func MakeGameServer(iterations, turns int) *GameServer {
	// embed BaseServer: maxTimeout = 10ms, maxThreads = 100
	serv := &GameServer{
		BaseServer: server.CreateBaseServer[infra.IGameAgent](iterations, turns, 10*time.Millisecond, 100),
		Env:        infra.NewEnvironment(),
	}

	// set GameRunner to bind RunTurn to BaseServer
	serv.SetGameRunner(serv)

	// Initialise socket
	if err := serv.InitSocket("127.0.0.1:5000"); err != nil {
		fmt.Fprintf(os.Stderr, "Socket init failed: %v\n", err)
	}

	return serv
}
