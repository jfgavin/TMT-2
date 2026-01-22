package config

type Config struct {
	Iterations     int
	Turns          int
	NumAgents      int
	StartingEnergy int
	GridSize       int
}

func NewConfig() Config {
	return Config{
		Iterations:     2,
		Turns:          5,
		NumAgents:      4,
		StartingEnergy: 5,
		GridSize:       32,
	}
}
