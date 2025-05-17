package kvNode

import (
	"errors"
	"github.com/Amirali-Amirifar/kv/internal/config"
)

type Service struct {
	config *config.KvNodeConfig
	store  *Storage
	state  NodeState
}

func NewKvNodeService(config *config.KvNodeConfig) *Service {
	// TODO add start function,
	// TODO call controller
	return &Service{
		config: config,
		store:  NewNodeStore(),
	}
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

	return nil
}

func (k *Service) Del(key string) error {
	k.store.Delete(key)
	return nil
}
