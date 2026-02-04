package env

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

type Cluster struct {
	id     uuid.UUID
	center Position
}

func (env *Environment) AddClusterResources(id uuid.UUID, amt int) bool {
	// Check cluster exists
	cluster, ok := env.clusters[id]
	if !ok {
		return false
	}

	cfg := env.cfg.Resources
	radius, lambda := float64(cfg.Radius), float64(cfg.Lambda)

	maxTerm := 1 - math.Exp(-radius/lambda)

	for amt > 0 {
		// Random angle
		theta := rand.Float64() * 2 * math.Pi

		// Random distance from centre
		u := rand.Float64()
		dist := -lambda * math.Log(1-u*maxTerm)

		// Clamping
		if dist > radius {
			dist = radius
		}

		// Final position of new resource
		x := float64(cluster.center.X) + dist*math.Cos(theta)
		y := float64(cluster.center.Y) + dist*math.Sin(theta)

		newPos := Position{
			X: int(math.Round(x)),
			Y: int(math.Round(y)),
		}

		// Modify tile
		if ok := env.ChangeResources(newPos, 1); ok {
			amt--
		}
	}
	return true
}

func (env *Environment) NewCluster() uuid.UUID {
	cfg := env.cfg.Resources

	cluster := Cluster{
		id:     uuid.New(),
		center: env.GetRandPosPadded(cfg.Radius),
	}

	env.clusters[cluster.id] = cluster
	return cluster.id
}

func (env *Environment) IntroduceResources() {
	cfg := env.cfg.Resources

	clusterIDs := make([]uuid.UUID, cfg.ClusterCount)

	for i := range cfg.ClusterCount {
		id := env.NewCluster()
		clusterIDs[i] = id
	}

	initResources := cfg.ResourceCount
	for initResources > 0 {
		id := clusterIDs[rand.Intn(len(clusterIDs))]
		if ok := env.AddClusterResources(id, 1); ok {
			initResources--
		}
	}
}
