package internal

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "ACTIVE"
	NodeStatusInactive NodeStatus = "INACTIVE"
	NodeStatusFailed   NodeStatus = "FAILED"
)

type StoreNodeType string

const (
	NodeTypeMaster   StoreNodeType = "MASTER"
	NodeTypeFollower StoreNodeType = "FOLLOWER"
)

type NodeType string

const (
	NodeTypeKvStore      NodeType = "KV_STORE"
	NodeTypeController   NodeType = "CONTROLLER"
	NodeTypeLoadBalancer NodeType = "LOAD_BALANCER"
	NodeTypeClient       NodeType = "CLIENT"
)
