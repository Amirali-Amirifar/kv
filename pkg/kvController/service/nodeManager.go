package service

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Amirali-Amirifar/kv/internal"
)

type NodeInfo struct {
	ID            int
	ShardKey      int
	Status        internal.NodeStatus
	Address       net.TCPAddr
	StoreNodeType internal.StoreNodeType
	LastChecked   time.Time
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

func (nm *NodeManager) initializeNodes() {
	numNodes := nm.partitions * nm.replicas
	nodes := make([]*NodeInfo, 0, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, &NodeInfo{
			ID:            i,
			ShardKey:      i % nm.partitions,
			Status:        internal.NodeStatusUnregistered,
			Address:       net.TCPAddr{},
			StoreNodeType: internal.NodeTypeUnknown,
		})
	}
	nm.Nodes = nodes
}

func (nm *NodeManager) RegisterNode(address string, port int) error {
	// See if there is an empty spot in the nodes list,
	// unregistered / failed nodes are empty spots
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
		if node.Address.IP.Equal(addr.IP) && node.Address.Port == addr.Port && node.Status == internal.NodeStatusActive {
			node.LastChecked = time.Now()
			return fmt.Errorf("node %s:%d is already registered", address, port)
		}
	}

	for _, node := range nm.Nodes {
		if node.Status == internal.NodeStatusUnregistered || node.Status == internal.NodeStatusFailed {
			node.Address = addr
			node.Status = internal.NodeStatusSyncing
			// get data and sync with master partition
			node.LastChecked = time.Now()
			return nil
		}
	}

	return fmt.Errorf("cannot register node at %s:%d: all cluster spots are full", address, port)
}

func (nm *NodeManager) GetNodeInfo(nodeID int) (NodeInfo, error) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if nodeID < 0 || nodeID >= len(nm.Nodes) {
		return NodeInfo{}, fmt.Errorf("invalid node ID: %d", nodeID)
	}
	return *nm.Nodes[nodeID], nil
}

func (nm *NodeManager) GetActiveNodes() []NodeInfo {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	active := []NodeInfo{}
	for _, node := range nm.Nodes {
		if node.Status == internal.NodeStatusActive {
			active = append(active, *node)
		}
	}
	return active
}
