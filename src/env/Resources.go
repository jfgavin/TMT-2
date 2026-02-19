package env

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

// === Cluster type & Methods ===

type Cluster struct {
	id        uuid.UUID
	resources map[Position]int
	center    Position
	radius    float64
	lambda    float64
}

func (clu *Cluster) AddResources(amt int) {
	maxTerm := 1 - math.Exp(-clu.radius/clu.lambda)

	for amt > 0 {
		// Random angle
		theta := rand.Float64() * 2 * math.Pi

		// Random distance from centre
		u := rand.Float64()
		dist := min(-clu.lambda*math.Log(1-u*maxTerm), clu.radius)

		// Final position of new resource
		x := float64(clu.center.X) + dist*math.Cos(theta)
		y := float64(clu.center.Y) + dist*math.Sin(theta)

		newPos := Position{
			X: int(math.Round(x)),
			Y: int(math.Round(y)),
		}

		clu.resources[newPos]++
		amt--
	}
}

func (clu *Cluster) SubResource(pos Position) bool {
	subbed := false
	amt, ok := clu.resources[pos]
	if !ok {
		return false
	}
	if amt > 0 {
		clu.resources[pos]--
		subbed = true
	}
	if amt <= 0 {
		delete(clu.resources, pos)
	}
	return subbed
}

// === Environment Top-level Resources Methods ===

func (env *Environment) NewCluster() *Cluster {
	cfg := env.cfg.Resources
	rad := float64(cfg.Radius)

	cluster := &Cluster{
		id:        uuid.New(),
		resources: make(map[Position]int),
		center:    env.GetRandPosPadded(cfg.Radius),
		radius:    rad,
		lambda:    rad * cfg.LambdaRatio,
	}

	env.clusters[cluster.id] = cluster
	return cluster
}

func (env *Environment) GetCluster(id uuid.UUID) (*Cluster, bool) {
	cluster, ok := env.clusters[id]
	if ok {
		return cluster, true
	}
	return nil, false
}

func (env *Environment) IntroduceResources() {
	cfg := env.cfg.Resources

	clus := make([]*Cluster, cfg.ClusterCount)

	for i := range cfg.ClusterCount {
		clus[i] = env.NewCluster()
	}

	initResources := cfg.ResourceCount
	for initResources > 0 {
		clu := clus[rand.Intn(cfg.ClusterCount)]
		clu.AddResources(1)
		initResources--
	}
}

func (env *Environment) GetResources() map[Position]int {
	res := make(map[Position]int)
	for _, clu := range env.clusters {
		for pos, amt := range clu.resources {
			res[pos] += amt
		}
	}
	return res
}

func (env *Environment) DrainResources(pos Position, amt int) bool {
	for amt > 0 {
		clus := make([]*Cluster, 0)

		for _, clu := range env.clusters {
			if res, ok := clu.resources[pos]; ok && res > 0 {
				clus = append(clus, clu)
			}
		}

		lc := len(clus)

		if lc == 0 {
			return false
		}

		clu := clus[rand.Intn(lc)]
		subbed := clu.SubResource(pos)
		if subbed {
			amt--
		}

	}

	return true
}
