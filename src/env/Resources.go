package env

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

// === Cluster type & Methods ===

type Cluster struct {
	env    IEnvironment
	id     uuid.UUID
	center Position
	radius float64
	lambda float64
}

// Method which takes function handle, and enacts it on every tile possibly part of the cluster
func (clu *Cluster) ForEachTile(fn func(tile *Tile, id uuid.UUID)) {
	rad, cen := int(clu.radius+1), clu.center
	gs := clu.env.GridSize()

	for y := cen.Y - rad; y < cen.Y+rad; y++ {
		for x := cen.X - rad; x < cen.X+rad; x++ {
			pos := Position{x, y}
			if !pos.IsBounded(gs) {
				continue
			}
			if tile, ok := clu.env.GetTile(Position{x, y}); ok {
				fn(tile, clu.id)
			}
		}
	}
}

func (clu *Cluster) AddResources(amt int) bool {
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

		// Modify tile
		if tile, ok := clu.env.GetTile(newPos); ok {
			tile.AddResources(clu.id, 1)
			amt--
		}
	}
	return true
}

func (clu *Cluster) DecayCluster() {
	clu.ForEachTile(func(tile *Tile, id uuid.UUID) {
		tile.SubResources(id, 1)
	})
}

func (clu *Cluster) GetClusterTotal() int {
	sum := 0
	clu.ForEachTile(func(tile *Tile, id uuid.UUID) {
		if val, ok := tile.GetContributions(id); ok {
			sum += val
		}
	})
	return sum
}

// === Environment Top-level Resources Methods ===

func (env *Environment) NewCluster() *Cluster {
	cfg := env.cfg.Resources
	rad := float64(cfg.Radius)

	cluster := &Cluster{
		env:    env,
		id:     uuid.New(),
		center: env.GetRandPosPadded(cfg.Radius),
		radius: rad,
		lambda: rad * cfg.LambdaRatio,
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

	clusters := make([]*Cluster, cfg.ClusterCount)

	for i := range cfg.ClusterCount {
		clu := env.NewCluster()
		clusters[i] = clu
	}

	initResources := cfg.ResourceCount
	for initResources > 0 {
		clu := clusters[rand.Intn(cfg.ClusterCount)]
		clu.AddResources(1)
		initResources--
	}
}
