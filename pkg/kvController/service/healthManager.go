package service

import (
	"fmt"
	"net"
	"time"

	"github.com/Amirali-Amirifar/kv/internal"
)

type HealthManager struct {
	nodeManager *NodeManager
	interval    time.Duration
	stopChan    chan struct{}
}

func NewHealthManager(nodeManager *NodeManager, interval time.Duration) *HealthManager {
	return &HealthManager{
		nodeManager: nodeManager,
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

func (hm *HealthManager) Start() {
	go hm.healthCheckLoop()
}

func (hm *HealthManager) Stop() {
	close(hm.stopChan)
}

func (hm *HealthManager) healthCheckLoop() {
	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-hm.stopChan:
			return
		case <-ticker.C:
			hm.checkNodes()
		}
	}
}

func (hm *HealthManager) checkNodes() {
	nodes := hm.nodeManager.GetActiveNodes()
	for _, node := range nodes {
		if err := hm.checkNode(node); err != nil {
			// If node is unresponsive, mark it as failed
			hm.nodeManager.mutex.Lock()
			hm.nodeManager.Nodes[node.ID].Status = internal.NodeStatusFailed
			hm.nodeManager.mutex.Unlock()
		}
	}
}

func (hm *HealthManager) checkNode(node NodeInfo) error {
	// Try to establish a TCP connection to the node
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", node.Address.IP.String(), node.Address.Port), 2*time.Second)
	if err != nil {
		return fmt.Errorf("node health check failed: %v", err)
	}
	defer conn.Close()
	// TODO: Implement actual health check protocol
	// For now, just establishing a connection is enough to consider the node healthy
	return nil
}
