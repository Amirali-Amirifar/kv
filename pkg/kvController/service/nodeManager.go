package service

import (
	"net"

	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/pkg/kvNode"
)

type NodeManager struct {
	ID            int
	ShardKey      int
	ShardWAL      *kvNode.WAL
	Status        internal.NodeStatus
	Address       net.TCPAddr
	StoreNodeType internal.StoreNodeType
}

func NewNodeManager(id int, shardKey int, shardWAL *kvNode.WAL) *NodeManager {
	nm := &NodeManager{
		ID:            id,
		ShardKey:      shardKey,
		ShardWAL:      shardWAL,
		Status:        internal.NodeStatusUnregistered,
		Address:       net.TCPAddr{},
		StoreNodeType: internal.NodeTypeUnknown,
	}
	return nm
}
