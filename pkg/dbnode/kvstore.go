package dbnode

import (
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal"
)

func (p *PartitionHandler) Set(key, value string) error {
	if p.Role != PartitionRoleMaster {
		return fmt.Errorf("partition %d is not master", p.ID)
	}
	if p.status != internal.NodeStatusActive {
		return fmt.Errorf("partition %d is not active", p.ID)
	}

	p.dataMu.Lock()
	defer p.dataMu.Unlock()
	p.Data[key] = value

	log := fmt.Sprintf("SET %s=%s", key, value)
	p.appendToWAL(log)
	p.replicateLog(log)

	return nil
}

func (p *PartitionHandler) Get(key string) (string, bool) {
	if p.status != internal.NodeStatusActive {
		return "", false
	}

	p.dataMu.RLock()
	defer p.dataMu.RUnlock()
	val, ok := p.Data[key]
	return val, ok
}

func (p *ReplicaHandler) Get(key string) (string, bool) {
	if p.status != internal.NodeStatusActive {
		return "", false
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.Data[key]
	return val, ok
}

func (p *PartitionHandler) Delete(key string) error {
	if p.Role != PartitionRoleMaster {
		return fmt.Errorf("partition %d is not master", p.ID)
	}

	p.dataMu.Lock()
	delete(p.Data, key)
	p.dataMu.Unlock()

	log := fmt.Sprintf("DELETE %s", key)
	p.appendToWAL(log)
	p.replicateLog(log)

	return nil
}

// appendToWAL safely appends an entry to WAL.
func (p *PartitionHandler) appendToWAL(entry string) {
	p.walMu.Lock()
	defer p.walMu.Unlock()
	p.WAL.Logs = append(p.WAL.Logs, entry)
}

// replicateLog sends the log entry to all replicas.
func (p *PartitionHandler) replicateLog(entry string) {
	for _, replica := range p.Replicas {
		_ = replica.ApplyLog(entry)
	}
}
