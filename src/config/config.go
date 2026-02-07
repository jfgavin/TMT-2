package config

type Config struct {
	Serv  ServConfig
	Agent AgentConfig
	Env   EnvConfig
}

type ServConfig struct {
	Iterations int
	Turns      int
	NumAgents  int
}

type AgentConfig struct {
	StartingEnergy int
	VisualRange    int
	ResourceYield  int
}

type EnvConfig struct {
	GridSize  int
	Resources ResourcesConfig
}

type ResourcesConfig struct {
	ResourceCount int
	ClusterCount  int
	Radius        int
	LambdaRatio   float64
}

func NewConfig() Config {
	return Config{
		Serv: ServConfig{
			Iterations: 2,
			Turns:      50,
			NumAgents:  4,
		},
		Agent: AgentConfig{
			StartingEnergy: 25,
			VisualRange:    20,
			ResourceYield:  3,
		},
		Env: EnvConfig{
			GridSize: 16,
			Resources: ResourcesConfig{
				ResourceCount: 300,
				ClusterCount:  3,
				Radius:        4,
				LambdaRatio:   0.5,
			},
		},
	}
}
