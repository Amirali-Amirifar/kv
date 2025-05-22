package service

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/pkg/kvController/interfaces"
)

type NodeManagerInterface interface {
	GetShardInfo(shardID int) (interface {
		GetMaster() interface {
			GetID() int
			GetAddress() (string, int)
			GetStatus() internal.NodeStatus
		}
		GetFollowers() []interface {
			GetID() int
			GetAddress() (string, int)
			GetStatus() internal.NodeStatus
		}
	}, bool)
	GetNodeInfo(nodeID int) (NodeInfo, error)
	RegisterNode(address string, port int) error
	UpdateShardMaster(shardID int, masterID int) error
}

type ShardInfo struct {
	ShardKey  int
	Master    *NodeInfo
	Followers []*NodeInfo
}

type NodeInfo struct {
	ID            int
	ShardKey      int
	Status        internal.NodeStatus
	Address       net.TCPAddr
	StoreNodeType internal.StoreNodeType
}

type NodeManager struct {
	replicas      int
	partitions    int
	Nodes         []*NodeInfo
	mutex         sync.Mutex
	ShardMap      map[int]*ShardInfo
	timeout       time.Duration
	healthManager *HealthManager
}

func NewNodeManager(partitions int, replicas int, cfg *config.KvControllerConfig) *NodeManager {
	nm := &NodeManager{
		replicas:   replicas,
		partitions: partitions,
		mutex:      sync.Mutex{},
		timeout:    time.Duration(cfg.Discovery.HeartbeatIntervalMs) * time.Millisecond,
	}
	nm.initializeNodes()
	return nm
}

func (nm *NodeManager) initializeNodes() {
	numNodes := nm.partitions * nm.replicas
	nodes := make([]*NodeInfo, 0, numNodes)

	// Step 1: Initialize nodes and assign shard keys
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

	nm.ShardMap = make(map[int]*ShardInfo)

	for _, node := range nm.Nodes {
		shardKey := node.ShardKey

		if shardInfo, exists := nm.ShardMap[shardKey]; !exists {
			// First node for this shard: make it the leader
			node.StoreNodeType = internal.NodeTypeMaster
			nm.ShardMap[shardKey] = &ShardInfo{
				ShardKey:  shardKey,
				Master:    node,
				Followers: []*NodeInfo{},
			}
		} else {
			// Next nodes are replicas
			node.StoreNodeType = internal.NodeTypeFollower
			shardInfo.Followers = append(shardInfo.Followers, node)
		}
	}
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
		if node.Address.IP.Equal(addr.IP) && node.Address.Port == addr.Port {
			if node.Status == internal.NodeStatusActive {
				return fmt.Errorf("node %s:%d is already registered.", address, port)
			}
			node.Status = internal.NodeStatusSyncing
			// TODO: Start syncing data from master.
			return nil
		}
	}
	for _, node := range nm.Nodes {
		if node.StoreNodeType == internal.NodeTypeFollower && node.Status == internal.NodeStatusFailed {
			node.Address = addr
			node.Status = internal.NodeStatusSyncing
			return nil
		}
	}
	for _, node := range nm.Nodes {
		if node.Status == internal.NodeStatusUnregistered {
			node.Address = addr
			node.Status = internal.NodeStatusSyncing
			return nil
		}
	}
	return fmt.Errorf("cannot register node at %s:%d: all cluster spots are full.", address, port)
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

func (nm *NodeManager) GetShardInfo(shardID int) (interfaces.ShardInterface, bool) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	shardInfo, exists := nm.ShardMap[shardID]
	if !exists {
		return nil, false
	}

	return shardInfo, true
}

func (nm *NodeManager) UpdateShardMaster(shardID int, masterID int) error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	shardInfo, exists := nm.ShardMap[shardID]
	if !exists {
		return fmt.Errorf("shard %d not found", shardID)
	}

	// Find the target node
	var targetNode *NodeInfo
	for _, follower := range shardInfo.Followers {
		if follower.ID == masterID {
			targetNode = follower
			break
		}
	}

	if targetNode == nil {
		return fmt.Errorf("node %d not found in shard %d", masterID, shardID)
	}

	// Update the shard's master
	oldMaster := shardInfo.Master
	shardInfo.Master = targetNode
	targetNode.StoreNodeType = internal.NodeTypeMaster

	// Remove the new leader from followers list
	newFollowers := make([]*NodeInfo, 0)
	for _, f := range shardInfo.Followers {
		if f.ID != targetNode.ID {
			newFollowers = append(newFollowers, f)
		}
	}
	shardInfo.Followers = newFollowers

	// Add the old master to followers list if it exists
	if oldMaster != nil {
		oldMaster.StoreNodeType = internal.NodeTypeFollower
		shardInfo.Followers = append(shardInfo.Followers, oldMaster)
	}

	return nil
}

// Add interface methods to ShardInfo
func (s *ShardInfo) GetMaster() interfaces.NodeInterface {
	return s.Master
}

func (s *ShardInfo) GetFollowers() []interfaces.NodeInterface {
	followers := make([]interfaces.NodeInterface, len(s.Followers))
	for i, f := range s.Followers {
		followers[i] = f
	}
	return followers
}

func (n *NodeInfo) GetID() int {
	return n.ID
}

func (n *NodeInfo) GetAddress() (string, int) {
	return n.Address.IP.String(), n.Address.Port
}

func (n *NodeInfo) GetStatus() internal.NodeStatus {
	return n.Status
}
