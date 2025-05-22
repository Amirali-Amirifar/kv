package kvNode

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/sirupsen/logrus"
)

type Service struct {
	config *config.KvNodeConfig
	state  NodeState
	store  *Storage
	wal    *WAL
	mu     sync.RWMutex
}

func NewKvNodeService(cfg *config.KvNodeConfig) *Service {
	svc := &Service{
		config: cfg,
		state: NodeState{
			IsMaster: false,
			ShardKey: 0,
		},
		store: NewNodeStore(),
		mu:    sync.RWMutex{},
	}

	if svc.state.IsMaster {
		svc.wal = NewWAL(svc.state.ShardKey)
	}

	return svc
}

func (k *Service) Start() error {
	panic("implement me")
}

func (k *Service) Get(key string) (string, error) {
	value, ok := k.store.Get(key)
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

func (k *Service) Set(key, value string) error {
	k.store.Set(key, value)
	if k.state.IsMaster {
		if k.wal != nil {
			k.wal.Append("SET", key, value)
		} else {
			return fmt.Errorf("WAL is nil.")
		}
	}
	return nil
}

func (k *Service) Del(key string) error {
	k.store.Delete(key)
	if k.state.IsMaster {
		if k.wal != nil {
			k.wal.Append("DELETE", key, "")
		} else {
			return fmt.Errorf("WAL is nil.")
		}
	}
	return nil
}

func (k *Service) GetLastSeq() int64 {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.wal == nil {
		return 0
	}
	return k.wal.GetLastSeq()
}

func (k *Service) BecomeLeader() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state.IsMaster {
		return errors.New("already a leader")
	}

	// Initialize WAL if not exists
	if k.wal == nil {
		k.wal = NewWAL(k.state.ShardKey)
	}

	// Update state
	k.state.IsMaster = true

	logrus.WithFields(logrus.Fields{
		"shardKey": k.state.ShardKey,
	}).Info("Node became leader")

	return nil
}

func (k *Service) GetWALSince(seq int64) []WALRecord {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.wal == nil {
		return nil
	}
	return k.wal.GetSince(seq)
}

func (k *Service) ApplyWALRecord(record WALRecord) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	switch record.Operation {
	case "SET":
		k.store.Set(record.Key, record.Value)
	case "DELETE":
		k.store.Delete(record.Key)
	default:
		return errors.New("unknown operation in WAL record")
	}

	return nil
}
