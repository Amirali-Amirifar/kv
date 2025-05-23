package types

type NodeStatus string

const (
	NodeStatusActive       NodeStatus = "ACTIVE"
	NodeStatusInactive     NodeStatus = "INACTIVE"
	NodeStatusFailed       NodeStatus = "FAILED"
	NodeStatusUnregistered NodeStatus = "UNREGISTERED"
	NodeStatusSyncing      NodeStatus = "SYNCING"
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

type WALRecord struct {
	Operation string `json:"operation"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Seq       int64  `json:"seq"`
}
