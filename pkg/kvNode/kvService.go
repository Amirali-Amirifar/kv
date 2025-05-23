package kvNode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/internal/types"
	"github.com/sirupsen/logrus"
)

type Service struct {
	config *config.KvNodeConfig
	state  NodeState
	store  *Storage
	wal    *WAL
	mu     sync.RWMutex
	client *http.Client
}

func NewKvNodeService(cfg *config.KvNodeConfig) *Service {
	timeout := time.Duration(cfg.HTTPTimeout) * time.Millisecond
	client := &http.Client{Timeout: timeout}

	svc := &Service{
		config: cfg,
		state: NodeState{
			IsMaster: false,
			ShardKey: 0,
		},
		store:  NewNodeStore(),
		mu:     sync.RWMutex{},
		client: client,
	}

	if svc.state.IsMaster {
		svc.wal = NewWAL(svc.state.ShardKey)
	}

	return svc
}

func (k *Service) Start() error {
	// Register with controller
	if err := k.RegisterWithController(); err != nil {
		return err
	}
	// Start WAL
	go k.syncWALPeriodically()
	return nil
}

func (k *Service) RegisterWithController() error {
	// Register with controller
	registerReq := struct {
		Ip   string `json:"ip"`
		Port int    `json:"port"`
	}{
		Ip:   k.config.Address.Host,
		Port: k.config.Address.Port,
	}

	body, err := json.Marshal(registerReq)
	if err != nil {
		return fmt.Errorf("failed to marshal register request: %v", err)
	}

	resp, err := k.client.Post(
		fmt.Sprintf("http://%s:%d/internal/nodes/register", k.config.Controller.Host, k.config.Controller.Port),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to register with controller: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register with controller: status %d", resp.StatusCode)
	}

	var nodeInfo struct {
		ID            int                 `json:"id"`
		ShardKey      int                 `json:"shard_key"`
		Status        types.NodeStatus    `json:"status"`
		StoreNodeType types.StoreNodeType `json:"store_node_type"`
		LeaderID      int                 `json:"leader_id"`
		LeaderAddress struct {
			IP   string `json:"ip"`
			Port int    `json:"port"`
		} `json:"leader_address,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&nodeInfo); err != nil {
		return fmt.Errorf("failed to decode node info: %v", err)
	}

	// Update node state
	k.state.NodeID = nodeInfo.ID
	k.state.ShardKey = nodeInfo.ShardKey
	k.state.LeaderID = nodeInfo.LeaderID

	// Update node type
	if nodeInfo.StoreNodeType == types.NodeTypeMaster {
		k.state.IsMaster = true
	} else {
		k.state.IsMaster = false
		k.state.MasterAddress = nodeInfo.LeaderAddress.IP
		k.state.MasterPort = nodeInfo.LeaderAddress.Port
	}

	return nil
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
	if k.state.IsMaster {
		if k.wal != nil {
			k.wal.Append("SET", key, value)
		} else {
			return fmt.Errorf("WAL is nil.")
		}
	}
	return nil
}

func (k *Service) Del(key string) error {
	k.store.Delete(key)
	if k.state.IsMaster {
		if k.wal != nil {
			k.wal.Append("DELETE", key, "")
		} else {
			return fmt.Errorf("WAL is nil.")
		}
	}
	return nil
}

func (k *Service) GetLastSeq() int64 {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.wal == nil {
		return 0
	}
	return k.wal.GetLastSeq()
}

func (k *Service) UpdateNodeState(state types.StoreNodeType, leaderID int) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if state == types.NodeTypeMaster {
		if k.state.IsMaster {
			return errors.New("already a leader")
		}
		if k.wal == nil {
			k.wal = NewWAL(k.state.ShardKey)
		}
		k.state.IsMaster = true

		logrus.WithFields(logrus.Fields{
			"shardKey": k.state.ShardKey,
		}).Info("Node became leader")
	} else if state == types.NodeTypeFollower {
		if !k.state.IsMaster {
			return errors.New("already a follower")
		}

		k.state.IsMaster = false
		k.state.LeaderID = leaderID
		resp, err := k.client.Get(fmt.Sprintf("http://%s:%d/node/%d", k.config.Controller.Host, k.config.Controller.Port, leaderID))
		if err != nil {
			return fmt.Errorf("failed to get master node info: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get master node info: status %d", resp.StatusCode)
		}

		var nodeInfo struct {
			Address struct {
				IP   string `json:"ip"`
				Port int    `json:"port"`
			} `json:"address"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&nodeInfo); err != nil {
			return fmt.Errorf("failed to decode master node info: %v", err)
		}

		// Update master address and port
		k.state.MasterAddress = nodeInfo.Address.IP
		k.state.MasterPort = nodeInfo.Address.Port

		logrus.WithFields(logrus.Fields{
			"shardKey": k.state.ShardKey,
			"leaderID": k.state.LeaderID,
			"master":   fmt.Sprintf("%s:%d", k.state.MasterAddress, k.state.MasterPort),
		}).Info("Node became follower")
	}

	return nil
}

func (k *Service) GetWALSince(seq int64) []WALRecord {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.wal == nil {
		return nil
	}
	return k.wal.GetSince(seq)
}

func (k *Service) ApplyWALRecord(record WALRecord) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	switch record.Operation {
	case "SET":
		k.store.Set(record.Key, record.Value)
	case "DELETE":
		k.store.Delete(record.Key)
	default:
		return errors.New("unknown operation in WAL record")
	}

	return nil
}

func (k *Service) syncWALPeriodically() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if k.state.IsMaster {
			if k.wal != nil {
				minSeq := k.wal.GetMinFollowerSeq()
				if minSeq > 0 {
					k.wal.ClearUntil(minSeq)
				}
			}
			return
		}

		// Get WAL entries from master
		resp, err := k.client.Get(fmt.Sprintf("http://%s:%d/wal/get-since/?since=%d", k.state.MasterAddress, k.state.MasterPort, k.state.LastWALSeq))
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"master": fmt.Sprintf("%s:%d", k.state.MasterAddress, k.state.MasterPort),
				"seq":    k.state.LastWALSeq,
			}).Error("Failed to fetch WAL entries from master")
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logrus.WithFields(logrus.Fields{
				"status": resp.StatusCode,
				"master": fmt.Sprintf("%s:%d", k.state.MasterAddress, k.state.MasterPort),
				"seq":    k.state.LastWALSeq,
			}).Error("Failed to fetch WAL entries from master")
			continue
		}

		var records []WALRecord
		if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"master": fmt.Sprintf("%s:%d", k.state.MasterAddress, k.state.MasterPort),
				"seq":    k.state.LastWALSeq,
			}).Error("Failed to decode WAL entries")
			continue
		}

		if len(records) == 0 {
			continue
		}

		sort.Slice(records, func(i, j int) bool {
			return records[i].Seq < records[j].Seq
		})

		for _, record := range records {
			if record.Seq <= k.state.LastWALSeq {
				continue
			}

			if err := k.ApplyWALRecord(record); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"seq":       record.Seq,
					"operation": record.Operation,
					"key":       record.Key,
				}).Error("Failed to apply WAL record")
				break
			}

			k.state.LastWALSeq = record.Seq

			// Notify master about our progress
			progressResp, err := k.client.Post(
				fmt.Sprintf("http://%s:%d/wal/progress", k.state.MasterAddress, k.state.MasterPort),
				"application/json",
				bytes.NewBufferString(fmt.Sprintf(`{"follower_id": %d, "seq": %d}`, k.state.NodeID, record.Seq)),
			)
			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"seq": record.Seq,
				}).Error("Failed to notify master about WAL progress")
				continue
			}
			progressResp.Body.Close()

			logrus.WithFields(logrus.Fields{
				"seq":       record.Seq,
				"operation": record.Operation,
				"key":       record.Key,
			}).Debug("Applied WAL record")
		}
	}
}

func (k *Service) UpdateFollowerProgress(followerID int, seq int64) {
	if k.wal != nil {
		k.wal.UpdateFollowerProgress(followerID, seq)
	}
}
