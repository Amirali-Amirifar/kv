package service

import (
	"fmt"
	"net"
	"sync"

	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/pkg/kvNode"
)

type SystemManager struct {
	replicas          int
	partitions        int
	PartitionManagers []*PartitionManager
	NodesManagers     []*NodeManager
	mutex             sync.Mutex
	lastWALSeq        int64
}

func NewSystemManager(partitions, replicas int) *SystemManager {
	sm := &SystemManager{
		replicas:   replicas,
		partitions: partitions,
	}
	sm.initializePartitions()
	sm.initializeNodes()
	return sm
}

func (sm *SystemManager) initializePartitions() {
	sm.PartitionManagers = make([]*PartitionManager, sm.partitions)
	for i := 0; i < sm.partitions; i++ {
		sm.PartitionManagers[i] = NewPartitionManager(i)
	}
}

func (sm *SystemManager) initializeNodes() {
	totalNodes := sm.partitions * sm.replicas
	sm.NodesManagers = make([]*NodeManager, totalNodes)
	for i := 0; i < totalNodes; i++ {
		PartitionID := i % sm.partitions
		WAL := sm.getWAL(PartitionID)
		sm.NodesManagers[i] = NewNodeManager(i, PartitionID, WAL)
	}
}

func (sm *SystemManager) RegisterNode(address string, port int) error {
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

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for _, node := range sm.NodesManagers {
		if node.Address.IP.Equal(addr.IP) && node.Address.Port == addr.Port && node.Status == internal.NodeStatusActive {
			return fmt.Errorf("node %s:%d is already registered", address, port)
		}
	}

	for _, node := range sm.NodesManagers {
		if node.Status == internal.NodeStatusUnregistered || node.Status == internal.NodeStatusFailed {
			node.Address = addr
			node.Status = internal.NodeStatusSyncing
			// TODO: Add data to syncing nodes
			return nil
		}
	}

	return fmt.Errorf("cannot register node at %s:%d: all cluster spots are full", address, port)
}

func (sm *SystemManager) GetNodeInfo(nodeID int) (NodeManager, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if nodeID < 0 || nodeID >= len(sm.NodesManagers) {
		return NodeManager{}, fmt.Errorf("invalid node ID: %d", nodeID)
	}
	return *sm.NodesManagers[nodeID], nil
}

func (sm *SystemManager) GetActiveNodes() []NodeManager {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	var active []NodeManager
	for _, node := range sm.NodesManagers {
		if node.Status == internal.NodeStatusActive {
			active = append(active, *node)
		}
	}
	return active
}

func (sm *SystemManager) getWAL(partitionID int) *kvNode.WAL {
	return &sm.PartitionManagers[partitionID].ParitionWAL
}
