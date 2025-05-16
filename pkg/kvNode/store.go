package kvNode

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var log = logrus.WithFields(logrus.Fields{
	"Service": "NodeStore",
})

type NodeStore struct {
	data map[string]string
	mu   *sync.RWMutex
}

func NewNodeStore() *NodeStore {
	return &NodeStore{
		data: make(map[string]string),
		mu:   &sync.RWMutex{},
	}
}

func (s *NodeStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	log.Printf("%+v\n", s.data)
}
func (s *NodeStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	log.Printf("%+v\n", s.data)
	return val, ok
}

func (s *NodeStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	log.Printf("%+v\n", s.data)
}
