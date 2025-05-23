package interfaces

import (
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
)

type NodeManagerInterface interface {
	GetShardInfo(shardID int) (*cluster.ShardInfo, bool)
	GetNodeInfo(nodeID int) (cluster.NodeInfo, error)
	RegisterNode(address string, port int) (*cluster.NodeInfo, error)
	UpdateShardMaster(shardID int, masterID int) error
	UpdateNodeStatus(nodeID int, status cluster.NodeStatus) error
}

type KvControllerInterface interface {
	RegisterNode(address string, port int) (*cluster.NodeInfo, error)
	ChangePartitionLeader(shardID int, nodeID int) error
	GetNodeManager() NodeManagerInterface
	UpdateNodeStatus(nodeID int, status cluster.NodeStatus) error
}
