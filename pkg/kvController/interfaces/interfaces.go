package interfaces

import "github.com/Amirali-Amirifar/kv/internal"

type NodeInterface interface {
	GetID() int
	GetAddress() (string, int)
	GetStatus() internal.NodeStatus
}

type ShardInterface interface {
	GetMaster() NodeInterface
	GetFollowers() []NodeInterface
}

type NodeManagerInterface interface {
	GetShardInfo(shardID int) (ShardInterface, bool)
}

type KvControllerInterface interface {
	RegisterNode(address string, port int) error
	ChangePartitionLeader(shardID int, nodeID int) error
	GetNodeManager() NodeManagerInterface
}
