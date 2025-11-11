package gameServer

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
	"unicode"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/jfgavin/TMT-2/agents"
	"github.com/jfgavin/TMT-2/infra"
)

type GameServer struct {
	// embed functionality from package...
	// ...and tell BaseServer we're using IGameAgent
	*server.BaseServer[infra.IGameAgent]

	gridWidth, gridHeight int
	grid                  map[infra.Position]infra.IGameAgent
	environment           *infra.Environment
	conn                  net.Conn
}

// RunTurn implementation - Move, show grid, and then Message
func (serv *GameServer) RunTurn(i, j int) {
	StreamGameIteration(serv, i, j)
	fmt.Printf("\n====================[ ITERATION %d, TURN %d ]====================\n", i, j)

	// Check if any hero has run out of energy - skip iteration if so
	if serv.ShouldSkipIteration() {
		fmt.Println("=== SKIPPING ITERATION: A hero has run out of energy ===")
		return
	}

	// Use intelligent movement strategy (heroes go to their own shops first, then exits)
	serv.SetMyIntelligentMovementDirections()
	serv.RunMovementTurn()

	// Check for loot shop interactions
	serv.CheckLootShopInteractions()

	// Check for exit conditions
	serv.CheckExitConditions()

	fmt.Println("===GRID===:")
	serv.PrintGrid()

	// Check win condition
	if serv.CheckWinCondition() {
		fmt.Println("\n=== VICTORY! All heroes have looted their shops and exited! ===")
		// Note: The game will continue until all iterations complete
		// In a full implementation, you might want to stop early
	}
}

// just implement GameRunner
func (serv *GameServer) RunStartOfIteration(int) {
	fmt.Print("\n====================[ ITERATION BEGINNING ]====================\n")
	// serv.SetRandomActions() // Use random strategy (original)
	serv.SetMyIntelligentActions() // Use my intelligent strategy

	for _, gc := range serv.GetAgentMap() {
		gc.ResetEnergy()
		// Reset exit status for new iteration
		gc.SetExited(false)
		// Note: Don't reset memory (KnownExits, KnownShops) - agents remember across iterations
	}
}

// just implement GameRunner
func (serv *GameServer) RunEndOfIteration(int) {
	fmt.Print("\n====================[ ITERATION OVER ]====================\n")
	// MVP Requirement: No communication during iteration
	// Communication happens between iterations
	fmt.Println("===MESSAGING (Between Iterations)===:")
	serv.RunMessagingTurn()
}

// output current grid of agents
func (serv *GameServer) PrintGrid() {
	fmt.Println("Legend: [d/e/t/o]ero, [X]xit, [S]hop (unlooted), [$]hop (looted), [.]Normal")
	fmt.Println("Current Grid State:")

	// Build a quick lookup from reported positions
	gridLookup := make(map[infra.Position]infra.IGameAgent)
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			gridLookup[gc.GetPos()] = gc
		}
	}

	for y := 0; y < serv.gridHeight; y++ {
		for x := 0; x < serv.gridWidth; x++ {
			pos := infra.Position{X: x, Y: y}
			if gc, ok := gridLookup[pos]; ok {
				// Hero is at this position - show hero name (lowercase to distinguish from Exit)
				heroChar := []rune(gc.GetName())[0]
				// Convert to lowercase for hero display
				heroChar = unicode.ToLower(heroChar)
				fmt.Printf("%c ", heroChar)
			} else {
				// No hero - show tile type
				tile := serv.environment.GetTile(pos)
				switch tile.Type {
				case infra.ExitTile:
					fmt.Print("X ") // Changed from "E" to "X" to avoid conflict
				case infra.LootShopTile:
					// Show 'S' for unlooted shop, '$' for looted shop
					if serv.environment.IsShopLooted(pos) {
						fmt.Print("$ ")
					} else {
						fmt.Print("S ")
					}
				default:
					fmt.Print(". ")
				}
			}
		}
		fmt.Println()
	}
}

// ========== RANDOM METHODS (Original Implementation) ==========
// SetRandomMovementDirections - Original random movement assignment
func (serv *GameServer) SetRandomMovementDirections() {
	movements := []infra.Direction{infra.North, infra.South, infra.West, infra.East}
	for _, gc := range serv.GetAgentMap() {
		// set random movement direction for this agent
		index := rand.Intn(len(movements))
		gc.SetMovementDir(movements[index])

		// delete this direction from the list, so it cannot be re-assigned this turn
		movements[index] = movements[len(movements)-1]
		movements = movements[:len(movements)-1]
	}
}

func (serv *GameServer) SetRandomActions() {
	// MVP: Actions are disabled (explore mall is not used in MVP)
	// Set all actions to NoAction
	for _, gc := range serv.GetAgentMap() {
		gc.SetAction(infra.NoAction)
	}
}

// ========== MY INTELLIGENT METHODS ==========
// SetMyIntelligentMovementDirections - My intelligent strategy for movement
// (Kept for future use - not used in MVP)
func (serv *GameServer) SetMyIntelligentMovementDirections() {
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			if uga, ok := gc.(*agents.UserGameAgent); ok {
				direction := uga.DecideMovementDirection()
				if direction != infra.NoDirection {
					gc.SetMovementDir(direction)
				} else {
					movements := []infra.Direction{infra.North, infra.South, infra.West, infra.East}
					gc.SetMovementDir(movements[rand.Intn(len(movements))])
				}
			} else {
				movements := []infra.Direction{infra.North, infra.South, infra.West, infra.East}
				gc.SetMovementDir(movements[rand.Intn(len(movements))])
			}
		}
	}
}

// SetMyIntelligentActions - My intelligent strategy for actions
// (Kept for future use - not used in MVP)
func (serv *GameServer) SetMyIntelligentActions() {
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			if uga, ok := gc.(*agents.UserGameAgent); ok {
				action := uga.DecideAction()
				gc.SetAction(action)
			} else {
				gc.SetAction(infra.NoAction)
			}
		}
	}
}

// make all agents move
func (serv *GameServer) RunMovementTurn() {
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			gc.DoMovement()
		}
	}
}

// make all agents perform their actions
// Note: MVP does not use actions (explore mall is disabled)
func (serv *GameServer) RunActionTurn() {
	// Actions are disabled in MVP - no action execution needed
}

// Check if heroes are on exit tiles and mark them as exited
func (serv *GameServer) CheckExitConditions() {
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() && serv.environment.IsExit(gc.GetPos()) {
			gc.SetExited(true)
			fmt.Printf("[%s] has reached an exit and exited!\n", gc.GetName())
		}
	}
}

// Check if heroes are on loot shop tiles and allow interaction
func (serv *GameServer) CheckLootShopInteractions() {
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			// Check if hero is at their own loot shop
			if serv.environment.IsMyLootShop(gc.GetPos(), gc.GetName()) {
				// Update hero's memory of their own shop
				if uga, ok := gc.(*agents.UserGameAgent); ok {
					uga.SetMyShop(gc.GetPos())

					// Mark shop as looted if not already looted
					if !uga.HasLootedMyShop {
						uga.HasLootedMyShop = true
						serv.environment.MarkShopAsLooted(gc.GetPos())
						fmt.Printf("[%s] has LOOTED their shop at (%d,%d)!\n",
							gc.GetName(), gc.GetPos().X, gc.GetPos().Y)
					}
				}
				// MVP: Loot shop interaction - restore energy (only at own shop)
				if gc.GetEnergy() < infra.StartingEnergy {
					gc.AddEnergy(1)
					fmt.Printf("[%s] is at their own loot shop and restored 1 energy! (Energy: %d/%d)\n",
						gc.GetName(), gc.GetEnergy(), infra.StartingEnergy)
				}
			} else if serv.environment.IsLootShop(gc.GetPos()) {
				// Hero is at someone else's shop - no interaction
				fmt.Printf("[%s] is at someone else's loot shop (no interaction)\n", gc.GetName())
			}
		}
	}
}

// make all agents message
func (serv *GameServer) RunMessagingTurn() {
	// Let agents signal they're ready for communication
	for _, gc := range serv.GetAgentMap() {
		gc.DoMessaging()
	}

	// Server-side coordination: Share information between all agents
	// This allows agents to work with other agents through server coordination
	serv.ShareInformationBetweenAgents()
}

// ========== WORK WITH OTHER AGENTS (Communication) ==========
// ShareInformationBetweenAgents shares information between all agents
// This uses GetAllAgents to work with other agents
func (serv *GameServer) ShareInformationBetweenAgents() {
	allAgents := serv.GetAgentMap()

	// Share information between all pairs of agents
	for _, agent1 := range allAgents {
		ugc1, ok1 := agent1.(*agents.UserGameAgent)
		if !ok1 {
			continue
		}

		for _, agent2 := range allAgents {
			if agent1 == agent2 {
				continue // Skip self
			}

			ugc2, ok2 := agent2.(*agents.UserGameAgent)
			if !ok2 {
				continue
			}

			// Share exits
			for _, exit := range ugc1.KnownExits {
				ugc2.AddKnownExit(exit)
			}

			// Share shops
			for _, shop := range ugc1.KnownShops {
				ugc2.AddKnownShop(shop)
			}
		}
	}
}

// Check if iteration should be skipped (any hero has 0 energy)
func (serv *GameServer) ShouldSkipIteration() bool {
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() && gc.GetEnergy() == 0 {
			return true
		}
	}
	return false
}

// Check win condition: all heroes have exited AND all shops have been looted
func (serv *GameServer) CheckWinCondition() bool {
	// Check if all heroes have exited
	for _, gc := range serv.GetAgentMap() {
		if !gc.HasExited() {
			return false
		}
	}
	// Check if all shops have been looted
	if !serv.environment.AllShopsLooted() {
		return false
	}
	return true
}

// Socket
func (serv *GameServer) InitSocket(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to Python: %w", err)
	} else {
		fmt.Printf("Socket successfully initialised at %s\n", address)
	}
	serv.conn = conn
	return nil
}

func (serv *GameServer) CloseSocket() {
	if serv.conn != nil {
		serv.conn.Close()
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
		gridWidth:  infra.StartingGridWidth,
		gridHeight: infra.StartingGridHeight,
		grid:       make(map[infra.Position]infra.IGameAgent),
	}

	// Create environment with mall tiles, exits, and loot shops
	// First collect hero names for shop assignment
	heroNames := []string{"Dwarf", "Elf", "Troll", "Orc"}
	serv.environment = infra.CreateEnvironment(serv.gridWidth, serv.gridHeight, heroNames)

	// set GameRunner to bind RunTurn to BaseServer
	serv.SetGameRunner(serv)

	// Initialise socket
	if err := serv.InitSocket("127.0.0.1:5000"); err != nil {
		fmt.Fprintf(os.Stderr, "Socket init failed: %v\n", err)
	}

	names := [4]string{"Dwarf", "Elf", "Troll", "Orc"}
	for i := 0; i < infra.NumAgents; i++ {
		// Place heroes at random positions, but not on exits
		var pos infra.Position
		for {
			pos = infra.Position{
				X: rand.Intn(serv.gridWidth),
				Y: rand.Intn(serv.gridHeight),
			}
			// Make sure starting position is not an exit
			if !serv.environment.IsExit(pos) {
				break
			}
		}

		var gc infra.IGameAgent = agents.GetUserGameAgent(serv, pos, names[i])

		// Set hero's own shop position and known exits (can only loot own shop)
		if uga, ok := gc.(*agents.UserGameAgent); ok {
			// Find the shop owned by this hero
			for _, shopPos := range serv.environment.Shops {
				if serv.environment.GetShopOwner(shopPos) == names[i] {
					uga.SetMyShop(shopPos)
					break
				}
			}
			// Heroes know all exits from the start (they can see the mall layout)
			for _, exitPos := range serv.environment.Exits {
				uga.AddKnownExit(exitPos)
			}
		}

		serv.grid[pos] = gc
		serv.AddAgent(gc)
	}

	return serv
}

// API Exposure

func (serv *GameServer) GridWidth() int {
	return serv.gridWidth
}

func (serv *GameServer) GridHeight() int {
	return serv.gridHeight
}

func (serv *GameServer) Environment() *infra.Environment {
	return serv.environment
}
