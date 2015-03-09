package strategy

import "github.com/elferink/swarm/cluster"

// WeightedNode represents a node in the cluster with a given weight, typically used for sorting
// purposes.
type MultiWeightedNode struct {
	Node cluster.Node
	// Weight is the inherent value of this node.
	CpuAndMemoryWeight int
	ContainerWeight int
}

type MultiWeightedNodeList []*MultiWeightedNode

func (n MultiWeightedNodeList) Len() int {
	return len(n)
}

func (n MultiWeightedNodeList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n MultiWeightedNodeList) Less(i, j int) bool {
	var (
		ip = n[i]
		jp = n[j]
	)

	if i.ContainerWeight != j.ContainerWeight {
		return i.ContainerWeight < j.ContainerWeight
	}
	return i.CpuAndMemoryWeight < j.CpuAndMemoryWeight
}
