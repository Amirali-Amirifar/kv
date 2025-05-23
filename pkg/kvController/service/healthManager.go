package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
	"net/http"
	"time"

	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/sirupsen/logrus"
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
			logrus.WithError(err).WithField("node", node.ID).Warn("Node health check failed")
			hm.handleNodeFailure(node)
		}
	}
}

func (hm *HealthManager) checkNode(node cluster.NodeInfo) error {
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

func (hm *HealthManager) handleNodeFailure(node cluster.NodeInfo) {
	hm.nodeManager.mutex.Lock()
	defer hm.nodeManager.mutex.Unlock()

	n := hm.nodeManager.Nodes[node.ID]
	if n != nil {
		n.Status = cluster.NodeStatusFailed

		// If this was a master node, we need to elect a new leader
		if n.StoreNodeType == cluster.NodeTypeMaster {
			hm.electNewLeader(n.ShardKey)
		}
	}
}

func (hm *HealthManager) electNewLeader(shardKey int) {
	shardInfo, exists := hm.nodeManager.ShardMap[shardKey]
	if !exists {
		logrus.WithField("shardKey", shardKey).Error("Shard not found during leader election")
		return
	}

	// Find the follower with the highest sequence number
	var newLeader *cluster.NodeInfo
	var highestSeq int64 = -1

	// Check all followers
	for _, follower := range shardInfo.Followers {
		if follower.Status != cluster.NodeStatusActive {
			continue
		}

		// Get the follower's last sequence number
		seq, err := hm.getNodeLastSeq(follower)
		if err != nil {
			logrus.WithError(err).WithField("node", follower.ID).Error("Failed to get node last sequence number")
			continue
		}
		if seq > highestSeq {
			highestSeq = seq
			newLeader = follower
		}
	}

	if newLeader == nil {
		logrus.WithField("shardKey", shardKey).Error("No suitable follower found for leader election, terminating the shard")
		// TODO: Remove the shard from the shard map
		return
	}

	// Update the shard's master
	shardInfo.Master = newLeader
	newLeader.StoreNodeType = cluster.NodeTypeMaster

	// Remove the new leader from followers list
	newFollowers := make([]*cluster.NodeInfo, 0)
	for _, f := range shardInfo.Followers {
		if f.ID != newLeader.ID {
			newFollowers = append(newFollowers, f)
		}
	}
	shardInfo.Followers = newFollowers

	// Notify the new leader
	if err := hm.notifyNewLeader(newLeader); err != nil {
		logrus.WithError(err).WithField("node", newLeader.ID).Error("Failed to notify new leader")
		return
	}

	// Notify followers about the leadership change
	if err := hm.notifyFollowers(shardInfo.Followers, newLeader); err != nil {
		logrus.WithError(err).Error("Failed to notify some followers about leader change")
	}

	logrus.WithFields(logrus.Fields{
		"shardKey":  shardKey,
		"newLeader": newLeader.ID,
	}).Info("New leader elected")
}

func (hm *HealthManager) updateNodeState(node *cluster.NodeInfo, state cluster.StoreNodeType, leaderID int) error {
	client := &http.Client{Timeout: hm.timeout}
	stateUpdate := struct {
		State    cluster.StoreNodeType `json:"state"`
		LeaderID int                   `json:"leader_id"`
	}{
		State:    state,
		LeaderID: leaderID,
	}

	body, err := json.Marshal(stateUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	resp, err := client.Post(
		fmt.Sprintf("http://%s:%d/update-state", node.Address.IP.String(), node.Address.Port),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}

func (hm *HealthManager) notifyNewLeader(node *cluster.NodeInfo) error {
	return hm.updateNodeState(node, cluster.NodeTypeMaster, node.ID)
}

func (hm *HealthManager) notifyFollowers(followers []*cluster.NodeInfo, newLeader *cluster.NodeInfo) error {
	var lastErr error
	for _, follower := range followers {
		if follower.ID == newLeader.ID {
			continue // Skip the new leader
		}
		if err := hm.updateNodeState(follower, cluster.NodeTypeFollower, newLeader.ID); err != nil {
			logrus.WithError(err).WithField("follower", follower.ID).Warn("Failed to notify follower about leader change")
			lastErr = err
		}
	}
	return lastErr
}

func (hm *HealthManager) getNodeLastSeq(node *cluster.NodeInfo) (int64, error) {
	client := &http.Client{Timeout: hm.timeout}
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/last-seq", node.Address.IP.String(), node.Address.Port))
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("node returned non-200 status: %d", resp.StatusCode)
	}

	var result struct {
		LastSeq int64 `json:"last_seq"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return -1, err
	}

	return result.LastSeq, nil
}
