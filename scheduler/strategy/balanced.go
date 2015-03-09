package strategy

import (
	"errors"
	"sort"

	"github.com/elferink/swarm/cluster"
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
	weightedNodes := MultiWeightedNodeList{}

	for _, node := range nodes {
		nodeMemory := node.TotalMemory()
		nodeCpus := node.TotalCpus()

		// Skip nodes that are smaller than the requested resources.
		if nodeMemory < int64(config.Memory) || nodeCpus < config.CpuShares {
			continue
		}

		var (
			cpuWeight    int = 100
			memoryWeight int = 100
		)

		if config.CpuShares > 0 {
			cpuWeight = (node.UsedCpus() + config.CpuShares) * 100 / nodeCpus
		}
		if config.Memory > 0 {
			memoryWeight = (node.UsedMemory() + config.Memory) * 100 / nodeMemory
		}

		if cpuWeight > 100 || memoryWeight > 100 {
			continue
		}

		var container = len(node.Containers())
		var cpuAndMem = cpuWeight + memoryWeight

		weightedNodes = append(weightedNodes, &MultiWeightedNode{Node: node, CpuAndMemoryWeight: cpuAndMem, ContainerWeight: container})
	}

	if len(weightedNodes) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	sort.Sort(weightedNodes)

	return weightedNodes[0].Node, nil
}
