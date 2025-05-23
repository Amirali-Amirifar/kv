package service

import (
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
	"net"
	"sync"
	"time"

	"github.com/Amirali-Amirifar/kv/internal/config"
)

type NodeManager struct {
	replicas      int
	partitions    int
	Nodes         []*cluster.NodeInfo
	mutex         sync.Mutex
	ShardMap      map[int]*cluster.ShardInfo
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
	nodes := make([]*cluster.NodeInfo, 0, numNodes)

	// Step 1: Initialize nodes and assign shard keys
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, &cluster.NodeInfo{
			ID:            i,
			ShardKey:      i % nm.partitions,
			Status:        cluster.NodeStatusUnregistered,
			Address:       net.TCPAddr{},
			StoreNodeType: cluster.NodeTypeUnknown,
		})
	}
	nm.Nodes = nodes

	nm.ShardMap = make(map[int]*cluster.ShardInfo)

	for _, node := range nm.Nodes {
		shardKey := node.ShardKey

		if shardInfo, exists := nm.ShardMap[shardKey]; !exists {
			// First node for this shard: make it the leader
			node.StoreNodeType = cluster.NodeTypeMaster
			nm.ShardMap[shardKey] = &cluster.ShardInfo{
				ShardKey:  shardKey,
				Master:    node,
				Followers: []*cluster.NodeInfo{},
			}
		} else {
			// Next nodes are replicas
			node.StoreNodeType = cluster.NodeTypeFollower
			shardInfo.Followers = append(shardInfo.Followers, node)
		}
	}
}

func (nm *NodeManager) RegisterNode(address string, port int) (*cluster.NodeInfo, error) {
	// See if there is an empty spot in the nodes list,
	// unregistered / failed nodes are empty spots
	ip := net.ParseIP(address)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", address)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", port)
	}
	addr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}

	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	for _, node := range nm.Nodes {
		if node.Address.IP.Equal(addr.IP) && node.Address.Port == addr.Port {
			if node.Status == cluster.NodeStatusActive {
				return nil, fmt.Errorf("node %s:%d is already registered.", address, port)
			}
			node.Status = cluster.NodeStatusSyncing
			// TODO: Start syncing data from master.
			return node, nil
		}
	}
	for _, node := range nm.Nodes {
		if node.StoreNodeType == cluster.NodeTypeFollower && node.Status == cluster.NodeStatusFailed {
			node.Address = addr
			node.Status = cluster.NodeStatusSyncing
			return node, nil
		}
	}
	for _, node := range nm.Nodes {
		if node.Status == cluster.NodeStatusUnregistered {
			node.Address = addr
			node.Status = cluster.NodeStatusSyncing
			return node, nil
		}
	}
	return nil, fmt.Errorf("cannot register node at %s:%d: all cluster spots are full", address, port)
}

func (nm *NodeManager) GetNodeInfo(nodeID int) (cluster.NodeInfo, error) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if nodeID < 0 || nodeID >= len(nm.Nodes) {
		return cluster.NodeInfo{}, fmt.Errorf("invalid node ID: %d", nodeID)
	}
	return *nm.Nodes[nodeID], nil
}

func (nm *NodeManager) GetActiveNodes() []cluster.NodeInfo {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	var active []cluster.NodeInfo
	for _, node := range nm.Nodes {
		if node.Status == cluster.NodeStatusActive {
			active = append(active, *node)
		}
	}
	return active
}

func (nm *NodeManager) GetShardInfo(shardID int) (*cluster.ShardInfo, bool) {
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
	var targetNode *cluster.NodeInfo
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
	targetNode.StoreNodeType = cluster.NodeTypeMaster

	// Remove the new leader from followers list
	newFollowers := make([]*cluster.NodeInfo, 0)
	for _, f := range shardInfo.Followers {
		if f.ID != targetNode.ID {
			newFollowers = append(newFollowers, f)
		}
	}
	shardInfo.Followers = newFollowers

	// Add the old master to followers list if it exists
	if oldMaster != nil {
		oldMaster.StoreNodeType = cluster.NodeTypeFollower
		shardInfo.Followers = append(shardInfo.Followers, oldMaster)
	}

	return nil
}
