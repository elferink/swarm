package strategy

import (
	"errors"
	"sort"

	"github.com/docker/swarm/cluster"
	"github.com/samalba/dockerclient"
)

var (
	ErrNoResourcesAvailable = errors.New("no resources available to schedule container")
)

type BalancedPlacementStrategy struct {
	overcommitRatio int64
}

func (p *BalancedPlacementStrategy) Initialize() error {
	return nil
}

func (p *BalancedPlacementStrategy) PlaceContainer(config *dockerclient.ContainerConfig, nodes []cluster.Node) (cluster.Node, error) {
	weightedNodes := weightedNodeList{}

	for _, node := range nodes {
		nodeMemory := node.TotalMemory()
		nodeCpus := node.TotalCpus()

		// Skip nodes that are smaller than the requested resources.
		if nodeMemory < int64(config.Memory) || nodeCpus < config.CpuShares {
			continue
		}

		var (
			cpuScore    int64 = 100
			memoryScore int64 = 100
		)

		if config.CpuShares > 0 {
			cpuScore = (node.UsedCpus() + config.CpuShares) * 100 / nodeCpus
		}
		if config.Memory > 0 {
			memoryScore = (node.UsedMemory() + config.Memory) * 100 / nodeMemory
		}

		if cpuScore > 100 || memoryScore > 100 {
			continue
		}

		var containerScore = int64(len(node.Containers())) + 1
		var total = cpuScore + memoryScore + containerScore

		weightedNodes = append(weightedNodes, &balancedScore{node: node, score: total})
	}

	if len(weightedNodes) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	// sort by highest weight
	sort.Sort(sort.Reverse(weightedNodes))

	return weightedNodes[0].Node, nil
}
