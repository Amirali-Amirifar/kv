package kvLoadbalancer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
	log "github.com/sirupsen/logrus"
	"hash/fnv"
	"net/http"
	"sync"

	"github.com/Amirali-Amirifar/kv/internal/config"
	apiTypes "github.com/Amirali-Amirifar/kv/internal/types/api"
	"github.com/Amirali-Amirifar/kv/pkg/kvLoadbalancer/api"
)

type LoadBalancerService struct {
	config     *config.KvLoadBalancerConfig
	shardNodes map[int]*ShardInfo
	client     *http.Client
	mu         sync.RWMutex
}

type ShardInfo struct {
	Master    *NodeInfo
	Followers []*NodeInfo
}

type NodeInfo struct {
	ID       int
	Address  string
	Port     int
	NodeType cluster.StoreNodeType
}

func NewLoadBalancerService(cfg *config.KvLoadBalancerConfig) *LoadBalancerService {

	svc := &LoadBalancerService{
		config:     cfg,
		shardNodes: make(map[int]*ShardInfo),
		client:     &http.Client{},
		mu:         sync.RWMutex{},
	}

	return svc
}

func (s *LoadBalancerService) Serve() {
	server := api.NewHTTPServer(s)
	err := server.Serve(s.config.Address.Port)
	if err != nil {
		panic(err)
		return
	}
}

func (s *LoadBalancerService) calculateShard(key string) int {
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		log.Fatalf("Error hashing key %v: %v", key, err)
		return 0
	}
	return int(h.Sum32()) % len(s.shardNodes)
}

func (s *LoadBalancerService) Get(key string) (string, error) {
	s.mu.RLock()
	shardID := s.calculateShard(key)
	shardInfo, exists := s.shardNodes[shardID]
	s.mu.RUnlock()

	if !exists || len(shardInfo.Followers) == 0 {
		return "", fmt.Errorf("no available nodes for shard %d", shardID)
	}

	// Round-robin between followers for read requests
	node := shardInfo.Master
	if len(shardInfo.Followers) > 0 {
		node = shardInfo.Followers[0] // TODO: Implement proper load balancing
	}

	req := apiTypes.GetRequest{Key: key}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/get", node.Address, node.Port),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("node returned status %d", resp.StatusCode)
	}

	var getResp apiTypes.GetResponse
	if err := json.NewDecoder(resp.Body).Decode(&getResp); err != nil {
		return "", err
	}

	return getResp.Value, nil
}

func (s *LoadBalancerService) Set(key, value string) error {
	s.mu.RLock()
	shardID := s.calculateShard(key)
	shardInfo, exists := s.shardNodes[shardID]
	s.mu.RUnlock()

	if !exists || shardInfo.Master == nil {
		return fmt.Errorf("no master node available for shard %d", shardID)
	}

	req := apiTypes.SetRequest{Key: key, Value: value}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/set", shardInfo.Master.Address, shardInfo.Master.Port),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("master node returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *LoadBalancerService) Del(key string) error {
	s.mu.RLock()
	shardID := s.calculateShard(key)
	shardInfo, exists := s.shardNodes[shardID]
	s.mu.RUnlock()

	if !exists || shardInfo.Master == nil {
		return fmt.Errorf("no master node available for shard %d", shardID)
	}

	req := apiTypes.DelRequest{Key: key}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/del", shardInfo.Master.Address, shardInfo.Master.Port),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("master node returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *LoadBalancerService) UpdateNodeData(shardID int, master *NodeInfo, followers []*NodeInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.shardNodes[shardID] = &ShardInfo{
		Master:    master,
		Followers: followers,
	}
}
