package dbnode

import (
	"github.com/Amirali-Amirifar/kv/internal"
	"sync"
	"time"
)

type Node struct {
	ID         int
	Address    string
	Status     internal.NodeStatus
	Partitions map[int]*PartitionHandler
	mu         sync.RWMutex
}

type PartitionRole string

const (
	PartitionRoleMaster  PartitionRole = "master"
	PartitionRoleReplica PartitionRole = "replica"
)

type PartitionHandler struct {
	ID        int
	Role      PartitionRole
	Replicas  map[int]*ReplicaHandler
	Data      map[string]string
	WAL       *WAL                //log of the partition
	Snapshots []*Snapshot         // snapshots of the partition
	dataMu    sync.RWMutex        // lock for data
	walMu     sync.Mutex          // lock for WAL
	status    internal.NodeStatus // Node Status
}

type ReplicaHandler struct {
	ID              int
	Data            map[string]string
	LastSyncedLog   int
	MasterPartition *PartitionHandler
	mu              sync.RWMutex
	status          internal.NodeStatus
	lastSyncTime    time.Time
}

type WAL struct {
	Logs []string
	mu   sync.Mutex
}

type Snapshot struct {
	ID        int
	Data      map[string]string
	Time      time.Time
	WALOffset int
}
