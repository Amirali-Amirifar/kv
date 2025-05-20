package kvNode

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{
	"Service": "Storage",
})

type Storage struct {
	data map[string]string
	mu   *sync.RWMutex
}

func NewNodeStore() *Storage {
	// TODO add start function,
	// TODO call controller
	return &Storage{
		data: make(map[string]string),
		mu:   &sync.RWMutex{},
	}
}

func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	log.Printf("%+v\n", s.data)
}
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	log.Printf("%+v\n", s.data)
	return val, ok
}

func (s *Storage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	log.Printf("%+v\n", s.data)
}
