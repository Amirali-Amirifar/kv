package kvLoadbalancer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
	log "github.com/sirupsen/logrus"
	"hash/fnv"
	"io"
	"net/http"
	"sync"

	"github.com/Amirali-Amirifar/kv/internal/config"
	apiTypes "github.com/Amirali-Amirifar/kv/internal/types/api"
	"github.com/Amirali-Amirifar/kv/pkg/kvLoadbalancer/api"
)

type LoadBalancerService struct {
	config     *config.KvLoadBalancerConfig
	shardNodes map[int]*cluster.ShardInfo
	client     *http.Client
	mu         sync.RWMutex
}

func NewLoadBalancerService(cfg *config.KvLoadBalancerConfig) *LoadBalancerService {

	svc := &LoadBalancerService{
		config:     cfg,
		shardNodes: make(map[int]*cluster.ShardInfo),
		client:     &http.Client{},
		mu:         sync.RWMutex{},
	}

	return svc
}

func (s *LoadBalancerService) Serve() {
	server := api.NewHTTPServer(s)
	go s.UpdateNodeData()
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
	//if len(shardInfo.Followers) > 0 {
	//	node = shardInfo.Followers[0] // TODO: Implement proper load balancing
	//}
	//
	req := apiTypes.GetRequest{Key: key}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/get", node.Address.IP, node.Address.Port),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("error getting key: %v", err)
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
		fmt.Sprintf("http://%s:%d/set", shardInfo.Master.Address.IP, shardInfo.Master.Address.Port),
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
		fmt.Sprintf("http://%s:%d/del", shardInfo.Master.Address, shardInfo.Master.Address.Port),
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

func (s *LoadBalancerService) UpdateNodeData() {
	type ClusterResponse struct {
		Shards map[string][]cluster.NodeInfo `json:"shards"`
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Make HTTP request to get cluster data
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/admin/cluster", s.config.Controller.Host, s.config.Controller.Port))
	if err != nil {
		log.Printf("Error getting cluster data: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d from cluster endpoint", resp.StatusCode)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	// Parse JSON response
	var clusterData ClusterResponse
	if err := json.Unmarshal(body, &clusterData); err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return
	}

	// Initialize shardNodes map if it doesn't exist
	if s.shardNodes == nil {
		s.shardNodes = make(map[int]*cluster.ShardInfo)
	}

	// Process each shard
	for shardKeyStr, nodes := range clusterData.Shards {
		if len(nodes) == 0 {
			continue
		}

		// Convert shard key from string to int (since JSON keys are strings)
		var shardKey int
		if _, err := fmt.Sscanf(shardKeyStr, "%d", &shardKey); err != nil {
			log.Printf("Error parsing shard key %s: %v", shardKeyStr, err)
			continue
		}

		var master *cluster.NodeInfo
		var followers []*cluster.NodeInfo

		// Categorize nodes based on node_type or leader_id
		for i := range nodes {
			node := &nodes[i]
			// Determine if node is master based on node_type or if it's the leader
			if node.StoreNodeType == cluster.NodeTypeMaster || node.ID == node.LeaderID {
				master = node
			} else {
				followers = append(followers, node)
			}
		}

		// Update shard information
		s.shardNodes[shardKey] = &cluster.ShardInfo{
			Master:    master,
			Followers: followers,
		}

		log.Printf("Updated shard %d: master=%v, followers=%d",
			shardKey,
			master != nil,
			len(followers))
	}

	log.Printf("Successfully updated cluster data for %d shards", len(s.shardNodes))
}
