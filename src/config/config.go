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
	Neurons        NeuronConfig
}

type NeuronConfig struct {
	TauRise  float64
	TauDecay float64
	Dt       float64
}

type EnvironmentConfig struct {
	GridSize      int
	GraveLifespan int
	Resources     ResourceConfig
}

type ResourceConfig struct {
	ResourceCount int
	ClusterCount  int
	Radius        int
	LambdaRatio   float64
}

func NewConfig() Config {
	return Config{
		Serv: ServerConfig{
			Iterations: 10,
			Turns:      50,
			NumAgents:  4,
		},
		Agent: AgentConfig{
			StartingEnergy: 25,
			VisualRange:    20,
			ResourceYield:  3,
			Neurons: NeuronConfig{
				TauRise:  10,
				TauDecay: 100,
				Dt:       0.1,
			},
		},
		Env: EnvironmentConfig{
			GridSize:      16,
			GraveLifespan: 5,
			Resources: ResourceConfig{
				ResourceCount: 300,
				ClusterCount:  3,
				Radius:        4,
				LambdaRatio:   0.5,
			},
		},
	}
}
