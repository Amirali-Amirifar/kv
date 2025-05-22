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

type NodeManager struct {
	replicas      int
	partitions    int
	Nodes         []*interfaces.NodeInfo
	mutex         sync.Mutex
	ShardMap      map[int]*interfaces.ShardInfo
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
	nodes := make([]*interfaces.NodeInfo, 0, numNodes)

	// Step 1: Initialize nodes and assign shard keys
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, &interfaces.NodeInfo{
			ID:            i,
			ShardKey:      i % nm.partitions,
			Status:        internal.NodeStatusUnregistered,
			Address:       net.TCPAddr{},
			StoreNodeType: internal.NodeTypeUnknown,
		})
	}
	nm.Nodes = nodes

	nm.ShardMap = make(map[int]*interfaces.ShardInfo)

	for _, node := range nm.Nodes {
		shardKey := node.ShardKey

		if shardInfo, exists := nm.ShardMap[shardKey]; !exists {
			// First node for this shard: make it the leader
			node.StoreNodeType = internal.NodeTypeMaster
			nm.ShardMap[shardKey] = &interfaces.ShardInfo{
				ShardKey:  shardKey,
				Master:    node,
				Followers: []*interfaces.NodeInfo{},
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

func (nm *NodeManager) GetNodeInfo(nodeID int) (interfaces.NodeInfo, error) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if nodeID < 0 || nodeID >= len(nm.Nodes) {
		return interfaces.NodeInfo{}, fmt.Errorf("invalid node ID: %d", nodeID)
	}
	return *nm.Nodes[nodeID], nil
}

func (nm *NodeManager) GetActiveNodes() []interfaces.NodeInfo {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	active := []interfaces.NodeInfo{}
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
	var targetNode *interfaces.NodeInfo
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
	newFollowers := make([]*interfaces.NodeInfo, 0)
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
