package interfaces

import (
	"net"

	"github.com/Amirali-Amirifar/kv/internal"
)

type NodeInfo struct {
	ID            int
	ShardKey      int
	Status        internal.NodeStatus
	Address       net.TCPAddr
	StoreNodeType internal.StoreNodeType
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

type ShardInfo struct {
	ShardKey  int
	Master    *NodeInfo
	Followers []*NodeInfo
}

func (s *ShardInfo) GetMaster() NodeInterface {
	return s.Master
}

func (s *ShardInfo) GetFollowers() []NodeInterface {
	followers := make([]NodeInterface, len(s.Followers))
	for i, f := range s.Followers {
		followers[i] = f
	}
	return followers
}

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
	GetNodeInfo(nodeID int) (NodeInfo, error)
	RegisterNode(address string, port int) error
	UpdateShardMaster(shardID int, masterID int) error
}

type KvControllerInterface interface {
	RegisterNode(address string, port int) error
	ChangePartitionLeader(shardID int, nodeID int) error
	GetNodeManager() NodeManagerInterface
}
