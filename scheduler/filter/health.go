package filter

import (
	"errors"

	"github.com/elferink/swarm/cluster"
	"github.com/samalba/dockerclient"
)

var (
	ErrNoHealthyNodeAvailable = errors.New("No healthy node available in the cluster")
)

// HealthFilter only schedules containers on healthy nodes.
type HealthFilter struct {
}

func (f *HealthFilter) Filter(_ *dockerclient.ContainerConfig, nodes []cluster.Node) ([]cluster.Node, error) {
	result := []cluster.Node{}
	for _, node := range nodes {
		if node.IsHealthy() {
			result = append(result, node)
		}
	}

	if len(result) == 0 {
		return nil, ErrNoHealthyNodeAvailable
	}

	return result, nil
}
