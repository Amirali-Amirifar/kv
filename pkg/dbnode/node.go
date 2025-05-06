package dbnode

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Node struct {
	ID         int
	Partitions map[int]*PartitionHandler
	mu         sync.RWMutex
}

type PartitionHandler struct {
	ID        int
	Replicas  map[int]*ReplicaHandler
	Data      map[string]string
	WAL       *WAL         //log of the partition
	Snapshots []*Snapshot  // snapshots of the partition
	dataMu    sync.RWMutex // lock for data
	walMu     sync.Mutex   // lock for WAL
}

type ReplicaHandler struct {
	ID              int
	Data            map[string]string
	LastSyncedLog   int
	MasterPartition *PartitionHandler
	mu              sync.RWMutex
}

type WAL struct {
	Logs []string
	mu   sync.Mutex
}

func (w *WAL) Append(entry string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Logs = append(w.Logs, entry)
}

type Snapshot struct {
	Data map[string]string
	Time time.Time
}

func (r *ReplicaHandler) SyncFromMaster() {
	master := r.MasterPartition
	master.walMu.Lock()
	defer master.walMu.Unlock()

	for i := r.LastSyncedLog + 1; i < len(master.WAL.Logs); i++ {
		entry := master.WAL.Logs[i]
		r.ApplyLog(entry)
		r.LastSyncedLog = i
	}
}

func (r *ReplicaHandler) ApplyLog(entry string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Example log format: "SET key=value", "GET key", "DELETE key"
	parts := strings.Fields(entry)
	if len(parts) < 2 {
		return fmt.Errorf("invalid log entry format: %s", entry)
	}

	operation := parts[0]
	keyValue := parts[1]

	switch operation {
	case "SET":
		// For SET operation, the entry will be in the format "key=value"
		keyValueParts := strings.SplitN(keyValue, "=", 2)
		if len(keyValueParts) != 2 {
			return fmt.Errorf("invalid SET log entry format: %s", entry)
		}
		key := keyValueParts[0]
		value := keyValueParts[1]
		r.Data[key] = value
		return nil

	case "DELETE":
		// For DELETE operation, the entry will be in the format "key"
		key := keyValue
		delete(r.Data, key)
		return nil

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
}
