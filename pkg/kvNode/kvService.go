package kvNode

import (
	"errors"

	"github.com/Amirali-Amirifar/kv/internal/config"
)

type Service struct {
	config *config.KvNodeConfig
	store  *Storage
	state  NodeState
	wal    *WAL
}

func NewKvNodeService(config *config.KvNodeConfig, state NodeState) *Service {
	svc := &Service{
		config: config,
		store:  NewNodeStore(),
		state:  state,
	}

	if state.IsMaster {
		svc.wal = NewWAL(state.ShardKey)
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

	if k.state.IsMaster && k.wal != nil {
		k.wal.Append("SET", key, value)
	}

	return nil
}

func (k *Service) Del(key string) error {
	k.store.Delete(key)

	if k.state.IsMaster && k.wal != nil {
		k.wal.Append("DELETE", key, "")
	}

	return nil
}

func (k *Service) GetWALSince(seq int64) []WALRecord {
	if k.wal == nil {
		return nil
	}
	return k.wal.GetSince(seq)
}
