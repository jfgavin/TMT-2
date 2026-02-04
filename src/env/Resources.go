package env

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

type Cluster struct {
	id      uuid.UUID
	centres []Position
}

func (env *Environment) GetClusterByID(id uuid.UUID) (Cluster, bool) {
	for _, clu := range env.clusters {
		if id == clu.id {
			return clu, true
		}
	}
	return Cluster{}, false
}

func (env *Environment) AddDistributedResource(pos Position, amt int) bool {
	if !pos.IsBounded(env.GridSize()) {
		return false
	}

	cfg := env.cfg.Resources
	radius, lambda := float64(cfg.Radius), float64(cfg.Lambda)
	maxTerm := 1 - math.Exp(-radius/lambda)

	// Random angle and distance
	theta := rand.Float64() * 2 * math.Pi
	u := rand.Float64()
	dist := -lambda * math.Log(1-u*maxTerm)
	if dist > radius {
		dist = radius
	}

	x := float64(pos.X) + dist*math.Cos(theta)
	y := float64(pos.Y) + dist*math.Sin(theta)

	newPos := Position{
		X: int(math.Round(x)),
		Y: int(math.Round(y)),
	}

	return env.ChangeResources(newPos, amt)
}

func (env *Environment) AddSubclusterResources(id uuid.UUID, cenInd, amt int) bool {
	cluster, ok := env.GetClusterByID(id)
	if !ok {
		return false
	}
	if cenInd >= len(cluster.centres) {
		return false
	}
	cenPos := cluster.centres[cenInd]
	return env.AddDistributedResource(cenPos, amt)
}

func (env *Environment) RandomlyAddResources(amt int) {
	for amt > 0 {
		cluster := env.clusters[rand.Intn(env.cfg.Resources.ClusterCount)]
		subInd := rand.Intn(len(cluster.centres))
		if ok := env.AddSubclusterResources(cluster.id, subInd, 1); ok {
			amt--
		}
	}
}

func (env *Environment) IntroduceResources() {
	cfg := env.cfg.Resources

	// Create clusters with sub-centres
	clusters := make([]Cluster, cfg.ClusterCount)
	for i := 0; i < cfg.ClusterCount; i++ {
		mainCentre := env.GetRandPosPadded(cfg.Radius)
		subCentres := make([]Position, cfg.SubClusterCount)

		for j := 0; j < cfg.SubClusterCount; j++ {
			offsetX := rand.Float64()*2*cfg.SubClusterOffset - cfg.SubClusterOffset
			offsetY := rand.Float64()*2*cfg.SubClusterOffset - cfg.SubClusterOffset

			subCentres[j] = Position{
				X: int(math.Round(float64(mainCentre.X) + offsetX)),
				Y: int(math.Round(float64(mainCentre.Y) + offsetY)),
			}
		}

		clusters[i] = Cluster{
			id:      uuid.New(),
			centres: subCentres,
		}
	}
	env.clusters = append(env.clusters, clusters...)
	env.RandomlyAddResources(cfg.ResourceCount)
}
