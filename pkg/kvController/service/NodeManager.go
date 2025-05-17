package service

import (
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal"
	"net"
	"sync"
)

type NodeInfo struct {
	*internal.KvStoreConfig
}
type NodeManager struct {
	replicas   int
	partitions int
	Nodes      []*NodeInfo
	mutex      sync.Mutex
}

func NewNodeManager(partitions int, replicas int) *NodeManager {
	nm := &NodeManager{
		replicas:   replicas,
		partitions: partitions,
		mutex:      sync.Mutex{},
	}
	nm.initializeNodes()

	return nm
}

func (nm *NodeManager) GetNodeInfo(name string) (NodeInfo, error) {
	return NodeInfo{}, nil
}

func (nm *NodeManager) initializeNodes() {
	numNodes := nm.partitions * nm.replicas
	nodes := make([]*NodeInfo, 0, numNodes)

	for i := range numNodes {
		nodes[i] = &NodeInfo{
			&internal.KvStoreConfig{
				ID:            i,
				Status:        internal.NodeStatusUnregistered,
				ShardKey:      i % nm.partitions,
				StoreNodeType: internal.NodeTypeUnknown,
				Address:       net.TCPAddr{},
			},
		}
	}
}

func (nm *NodeManager) RegisterNode(address string, port int) error {
	// See if there is an empty spot in the nodes list,
	// unregistered nodes are empty spots
	// See if there is any unregistered nodes
	ip := net.ParseIP(address)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", address)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %d", port)
	}
	addr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}

	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	for _, node := range nm.Nodes {
		if node.Status == internal.NodeStatusUnregistered {
			node.Address = addr
			node.Status = internal.NodeStatusActive
			return nil
		}
	}

	// See if there are any failed nodes
	for _, node := range nm.Nodes {
		if node.Status == internal.NodeStatusFailed {
			node.Status = internal.NodeStatusActive
			node.Address = addr
			return nil
		}
	}

	return fmt.Errorf("cannot register node at %s:%d: all cluster spots are full", address, port)
}
