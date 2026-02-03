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
	VisualRange    int
	ResourceYield  int
}

type EnvironmentConfig struct {
	GridSize  int
	Resources ResourceConfig
}

type ResourceConfig struct {
	ResourceCount int
	ClusterCount  int
	Radius        int
	Lambda        int
}

func NewConfig() Config {
	return Config{
		Serv: ServerConfig{
			Iterations: 2,
			Turns:      50,
			NumAgents:  4,
		},
		Agent: AgentConfig{
			StartingEnergy: 25,
			VisualRange:    20,
			ResourceYield:  3,
		},
		Env: EnvironmentConfig{
			GridSize: 64,
			Resources: ResourceConfig{
				ResourceCount: 100,
				ClusterCount:  4,
				Radius:        5,
				Lambda:        4,
			},
		},
	}
}
