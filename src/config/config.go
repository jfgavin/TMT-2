package config

type Config struct {
	Serv  ServerConfig
	Agent AgentConfig
	Env   EnvironmentConfig
}

type ServerConfig struct {
	Iterations int
	Turns      int
	NumAgents  int
}

type AgentConfig struct {
	StartingEnergy int
}

type EnvironmentConfig struct {
	GridSize int
}

func NewConfig() Config {
	return Config{
		Serv: ServerConfig{
			Iterations: 2,
			Turns:      5,
			NumAgents:  4,
		},
		Agent: AgentConfig{
			StartingEnergy: 5,
		},
		Env: EnvironmentConfig{
			GridSize: 64,
		},
	}
}
