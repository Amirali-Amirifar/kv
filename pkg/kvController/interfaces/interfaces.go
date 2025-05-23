package interfaces

import (
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
)

type NodeManagerInterface interface {
	GetShardInfo(shardID int) (cluster.ShardInfo, bool)
	GetNodeInfo(nodeID int) (cluster.NodeInfo, error)
	RegisterNode(address string, port int) error
	UpdateShardMaster(shardID int, masterID int) error
}

type KvControllerInterface interface {
	RegisterNode(address string, port int) error
	ChangePartitionLeader(shardID int, nodeID int) error
	GetNodeManager() NodeManagerInterface
}
