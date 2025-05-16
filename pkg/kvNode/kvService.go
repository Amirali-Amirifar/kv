package kvNode

import (
	"errors"
)

type Service struct {
	config *Config
	store  *NodeStore
}

func NewKvNodeService(config *Config) *Service {
	return &Service{
		config: config,
		store:  NewNodeStore(),
	}
}

func (k Service) Get(key string) (string, error) {
	value, ok := k.store.Get(key)
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

func (k Service) Set(key, value string) error {
	k.store.Set(key, value)

	return nil
}

func (k Service) Del(key string) error {
	k.store.Delete(key)

	return nil
}
