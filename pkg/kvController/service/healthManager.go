package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/internal/config"
)

type HealthManager struct {
	SystemManager *SystemManager
	interval      time.Duration
	timeout       time.Duration
	stopChan      chan struct{}
}

func NewHealthManager(systemManager *SystemManager, cfg *config.KvControllerConfig) *HealthManager {
	return &HealthManager{
		SystemManager: systemManager,
		interval:      time.Duration(cfg.Discovery.HeartbeatIntervalMs) * time.Millisecond,
		timeout:       time.Duration(cfg.Discovery.FailureTimeoutMs) * time.Millisecond,
		stopChan:      make(chan struct{}),
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

func (hm *HealthManager) checkNodes() map[int]internal.NodeStatus {
	nodes := hm.SystemManager.NodesManagers
	nodeStatus := make(map[int]internal.NodeStatus)

	for _, node := range nodes {
		nodeStatus[node.ID] = node.Status
	}
	return nodeStatus
}

func (hm *HealthManager) checkNode(node NodeManager) error {
	client := &http.Client{Timeout: hm.timeout}
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/health", node.Address.IP.String(), node.Address.Port))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}
