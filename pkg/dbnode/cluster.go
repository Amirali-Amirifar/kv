package dbnode

import (
	"fmt"
	"log"
	"strings"
)

func (r *ReplicaHandler) SyncFromMaster() {
	master := r.MasterPartition
	master.walMu.Lock()
	defer master.walMu.Unlock()

	for i := r.LastSyncedLog + 1; i < len(master.WAL.Logs); i++ {
		entry := master.WAL.Logs[i]
		if err := r.ApplyLog(entry); err != nil {
			log.Printf("replica %d failed to apply log at index %d: %v", r.ID, i, err)
			continue
		}
		r.LastSyncedLog = i
	}
}

func (r *ReplicaHandler) ApplyLog(entry string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	parts := strings.Fields(entry)
	if len(parts) < 2 {
		return fmt.Errorf("invalid log format: %s", entry)
	}

	op, payload := parts[0], parts[1]
	switch op {
	case "SET":
		return r.applySet(payload)
	case "DELETE":
		return r.applyDelete(payload)
	default:
		return fmt.Errorf("unsupported operation: %s", op)
	}
}

func (r *ReplicaHandler) applySet(kv string) error {
	parts := strings.SplitN(kv, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid SET format: %s", kv)
	}
	r.Data[parts[0]] = parts[1]
	return nil
}

func (r *ReplicaHandler) applyDelete(key string) error {
	delete(r.Data, key)
	return nil
}
