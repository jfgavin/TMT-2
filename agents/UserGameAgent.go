package agents

import (
	"fmt"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/jfgavin/TMT-2/infra"
)

// third tier of composition - embed BaseGameAgent..
// ... and add 'user specific' fields
type UserGameAgent struct {
	*infra.BaseGameAgent
	Name        string
	MovementDir infra.Direction
	Action      infra.Action

	// Memory and strategy
	KnownExits      []infra.Position // Known exit positions
	KnownShops      []infra.Position // Known shop positions (all shops)
	MyShop          *infra.Position  // My own shop position (can only loot this one)
	HasLootedMyShop bool             // Track if hero has looted their own shop
	TargetExit      *infra.Position  // Current target exit
	TargetShop      *infra.Position  // Current target shop (if energy low)
}

func (uga *UserGameAgent) GetName() string {
	return uga.Name
}

func (uga *UserGameAgent) SetAction(action infra.Action) {
	uga.Action = action
}

func (uga *UserGameAgent) GetAction() infra.Action {
	return uga.Action
}

func (uga *UserGameAgent) SetMovementDir(direction infra.Direction) {
	uga.MovementDir = direction
}

func (uga *UserGameAgent) GetMovementDir() infra.Direction {
	return uga.MovementDir
}

func (uga *UserGameAgent) DoMovement() {
	// Skip if already exited
	if uga.HasExited() {
		return
	}

	delta, ok := infra.Movement[uga.MovementDir]
	if !ok {
		return
	}

	if uga.Energy > 0 {
		uga.Pos.X += delta.X
		uga.Pos.Y += delta.Y

		// wrap around grid (Pac-Man style)
		uga.Pos.X = (uga.Pos.X + infra.StartingGridWidth) % infra.StartingGridWidth
		uga.Pos.Y = (uga.Pos.Y + infra.StartingGridHeight) % infra.StartingGridHeight

		uga.Energy--
	}
}

func (uga *UserGameAgent) DoAction() {
	// MVP: Actions are disabled (explore mall is not used)
	// No action execution needed
}

// DecideMovementDirection decides the best direction to move based on strategy
func (uga *UserGameAgent) DecideMovementDirection() infra.Direction {
	if uga.HasExited() {
		return infra.NoDirection
	}

	// Strategy: Prioritize reaching target
	var targetPos *infra.Position

	// Priority 1: If haven't looted my shop yet, go to my shop first (regardless of energy)
	if !uga.HasLootedMyShop && uga.MyShop != nil {
		targetPos = uga.MyShop
	}

	// Priority 2: If already looted my shop, go to exit (this is the main goal after looting)
	if targetPos == nil && uga.HasLootedMyShop {
		if len(uga.KnownExits) > 0 {
			nearestExit := uga.findNearestPosition(uga.KnownExits)
			if nearestExit != nil {
				targetPos = nearestExit
			}
		}
	}

	// Priority 3: If no specific target but we know exits, try to reach nearest exit
	if targetPos == nil && len(uga.KnownExits) > 0 {
		nearestExit := uga.findNearestPosition(uga.KnownExits)
		if nearestExit != nil {
			targetPos = nearestExit
		}
	}

	// Priority 4: Fallback - if we don't know our shop but know other shops, try any shop
	if targetPos == nil && uga.MyShop == nil && len(uga.KnownShops) > 0 {
		nearestShop := uga.findNearestPosition(uga.KnownShops)
		if nearestShop != nil {
			targetPos = nearestShop
		}
	}

	// If we have a target, move towards it
	if targetPos != nil {
		return uga.directionTowards(targetPos)
	}

	// For MVP, just move in a consistent pattern (North, then East, then South, then West)
	// This will be replaced by random if server needs fallback
	return infra.NoDirection // Server will use random fallback
}

// DecideAction decides what action to perform
func (uga *UserGameAgent) DecideAction() infra.Action {
	if uga.HasExited() || uga.Energy <= 0 {
		return infra.NoAction
	}

	// If we don't know any exits or shops, explore
	if len(uga.KnownExits) == 0 && len(uga.KnownShops) == 0 {
		if uga.Energy >= infra.ActionCost[infra.ExploreMall] {
			return infra.ExploreMall
		}
	}

	// If we know exits but energy is low, explore more to find shops
	if len(uga.KnownExits) > 0 && uga.Energy <= 2 && len(uga.KnownShops) == 0 {
		if uga.Energy >= infra.ActionCost[infra.ExploreMall] {
			return infra.ExploreMall
		}
	}

	// If we have target and enough energy, don't explore (save energy)
	if uga.TargetExit != nil || uga.TargetShop != nil {
		return infra.NoAction
	}

	// Default: explore if we have energy
	if uga.Energy >= infra.ActionCost[infra.ExploreMall] {
		return infra.ExploreMall
	}

	return infra.NoAction
}

// Helper: Find nearest position from a list
func (uga *UserGameAgent) findNearestPosition(positions []infra.Position) *infra.Position {
	if len(positions) == 0 {
		return nil
	}

	nearest := positions[0]
	minDist := uga.distance(uga.Pos, nearest)

	for _, pos := range positions[1:] {
		dist := uga.distance(uga.Pos, pos)
		if dist < minDist {
			minDist = dist
			nearest = pos
		}
	}

	return &nearest
}

// Helper: Calculate Manhattan distance
func (uga *UserGameAgent) distance(pos1, pos2 infra.Position) int {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// Helper: Get direction towards target
func (uga *UserGameAgent) directionTowards(target *infra.Position) infra.Direction {
	if target == nil {
		return infra.NoDirection
	}

	dx := target.X - uga.Pos.X
	dy := target.Y - uga.Pos.Y

	// Handle wrap-around (Pac-Man style) - choose shortest path
	// Calculate both direct and wrap-around distances, choose shorter
	dxDirect := abs(dx)
	dxWrap := abs(dx + infra.StartingGridWidth)
	if abs(dx-infra.StartingGridWidth) < dxWrap {
		dxWrap = abs(dx - infra.StartingGridWidth)
	}
	if dxWrap < dxDirect {
		if dx > 0 {
			dx -= infra.StartingGridWidth
		} else {
			dx += infra.StartingGridWidth
		}
	}

	dyDirect := abs(dy)
	dyWrap := abs(dy + infra.StartingGridHeight)
	if abs(dy-infra.StartingGridHeight) < dyWrap {
		dyWrap = abs(dy - infra.StartingGridHeight)
	}
	if dyWrap < dyDirect {
		if dy > 0 {
			dy -= infra.StartingGridHeight
		} else {
			dy += infra.StartingGridHeight
		}
	}

	// Prefer horizontal or vertical movement
	if abs(dx) > abs(dy) {
		if dx > 0 {
			return infra.East
		}
		return infra.West
	} else {
		if dy > 0 {
			return infra.South
		}
		return infra.North
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AddKnownExit adds an exit to memory
func (uga *UserGameAgent) AddKnownExit(pos infra.Position) {
	// Check if already known
	for _, known := range uga.KnownExits {
		if known.X == pos.X && known.Y == pos.Y {
			return
		}
	}
	uga.KnownExits = append(uga.KnownExits, pos)
	// Update target
	uga.TargetExit = uga.findNearestPosition(uga.KnownExits)
}

// AddKnownShop adds a shop to memory
func (uga *UserGameAgent) AddKnownShop(pos infra.Position) {
	// Check if already known
	for _, known := range uga.KnownShops {
		if known.X == pos.X && known.Y == pos.Y {
			return
		}
	}
	uga.KnownShops = append(uga.KnownShops, pos)
	// Update target
	uga.TargetShop = uga.findNearestPosition(uga.KnownShops)
}

// SetMyShop sets the hero's own shop position (can only loot this one)
func (uga *UserGameAgent) SetMyShop(pos infra.Position) {
	uga.MyShop = &pos
	// Also add to KnownShops for general knowledge
	uga.AddKnownShop(pos)
}

// Return if hero has looted their shop
func (uga *UserGameAgent) HasLoot() bool {
	return uga.HasLootedMyShop
}

// ProcessDiscovery processes discoveries from exploration
func (uga *UserGameAgent) ProcessDiscovery(discoveries []string) {
	for _, discovery := range discoveries {
		// Parse discovery string (format: "Exit at (x,y)" or "Loot Shop at (x,y)")
		var pos infra.Position
		var n int
		if n, _ = fmt.Sscanf(discovery, "Exit at (%d,%d)", &pos.X, &pos.Y); n == 2 {
			uga.AddKnownExit(pos)
		} else if n, _ = fmt.Sscanf(discovery, "Loot Shop at (%d,%d)", &pos.X, &pos.Y); n == 2 {
			uga.AddKnownShop(pos)
		}
	}
}

// DoMessaging overrides base implementation to share information with other agents
// Note: Communication with other agents is handled by server in RunMessagingTurn
// The server's ShareInformationBetweenAgents() method coordinates information sharing
func (uga *UserGameAgent) DoMessaging() {
	// Agent signals it's ready for communication
	// Server will coordinate information sharing between all agents
	// This allows agents to work with other agents through server coordination
	uga.SignalMessagingComplete()
}

// constructor for UserGameAgent
func GetUserGameAgent(funcs agent.IExposedServerFunctions[infra.IGameAgent], pos infra.Position, name string) *UserGameAgent {
	return &UserGameAgent{
		BaseGameAgent:   infra.GetBaseGameAgent(funcs, pos),
		Name:            name,
		MovementDir:     infra.NoDirection,
		Action:          infra.NoAction,
		KnownExits:      []infra.Position{},
		KnownShops:      []infra.Position{},
		MyShop:          nil,
		HasLootedMyShop: false,
		TargetExit:      nil,
		TargetShop:      nil,
	}
}
