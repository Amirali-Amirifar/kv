package service

import (
	"sync"

	"github.com/Amirali-Amirifar/kv/pkg/kvNode"
)

type PartitionManager struct {
	ID          int
	mutex       sync.Mutex
	ParitionWAL kvNode.WAL
}

func NewPartitionManager(ID int) *PartitionManager {
	pm := &PartitionManager{
		ID:          ID,
		mutex:       sync.Mutex{},
		ParitionWAL: *kvNode.NewWAL(ID),
	}
	return pm
}
