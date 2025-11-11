package infra

type Action int

const (
	ExploreMall Action = iota
	UseEscalator
	Warp
	ActivateTimer
	NoAction
)

var ActionCost = map[Action]int{
	ExploreMall:   1,
	UseEscalator:  2,
	Warp:          3,
	ActivateTimer: 2,
	NoAction:      0,
}

var ActionName = map[Action]string{
	ExploreMall:   "Explore Mall",
	UseEscalator:  "Use Escalator",
	Warp:          "Warp \t",
	ActivateTimer: "Activate Timer",
	NoAction:      "No Action",
}
