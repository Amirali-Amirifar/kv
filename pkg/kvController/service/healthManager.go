package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/internal/config"
)

type HealthManager struct {
	nodeManager *NodeManager
	interval    time.Duration
	timeout     time.Duration
	stopChan    chan struct{}
}

func NewHealthManager(nodeManager *NodeManager, cfg *config.KvControllerConfig) *HealthManager {
	return &HealthManager{
		nodeManager: nodeManager,
		interval:    time.Duration(cfg.Discovery.HeartbeatIntervalMs) * time.Millisecond,
		timeout:     time.Duration(cfg.Discovery.FailureTimeoutMs) * time.Millisecond,
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
			defer hm.nodeManager.mutex.Unlock()
			hm.nodeManager.Nodes[node.ID].Status = internal.NodeStatusFailed
		}
	}
}

func (hm *HealthManager) checkNode(node NodeInfo) error {
	client := &http.Client{Timeout: hm.timeout}
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/health", node.Address.IP.String(), node.Address.Port))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node returned non-200 status: %d", resp.StatusCode)
	}
	// for debugging
	node.LastChecked = time.Now()
	return nil
}
