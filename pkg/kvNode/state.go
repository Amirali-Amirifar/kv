package kvNode

// NodeState contains configs, and metadata fetched from controller
type NodeState struct {
	IsMaster      bool
	ShardKey      int
	LastWALSeq    int64
	LeaderID      int
	MasterAddress string
	MasterPort    int
	NodeID        int
}
