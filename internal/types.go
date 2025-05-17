package internal

import "net"

type NodeStatus string

const (
	NodeStatusActive       NodeStatus = "ACTIVE"
	NodeStatusInactive     NodeStatus = "INACTIVE"
	NodeStatusFailed       NodeStatus = "FAILED"
	NodeStatusUnregistered NodeStatus = "UNREGISTERED"
)

type StoreNodeType string

const (
	NodeTypeMaster   StoreNodeType = "MASTER"
	NodeTypeFollower StoreNodeType = "FOLLOWER"
	NodeTypeUnknown  StoreNodeType = "UNKNOWN"
)

type NodeType string

const (
	NodeTypeKvStore      NodeType = "KV_STORE"
	NodeTypeController   NodeType = "CONTROLLER"
	NodeTypeLoadBalancer NodeType = "LOAD_BALANCER"
	NodeTypeClient       NodeType = "CLIENT"
)

type KvStoreConfig struct {
	ID            int
	ShardKey      int
	Status        NodeStatus
	Address       net.TCPAddr
	StoreNodeType StoreNodeType
}
