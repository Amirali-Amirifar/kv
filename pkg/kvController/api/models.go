package api

import (
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
)

type NodeRegisterHandlerRequest struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type NodeRegisterHandlerResponse struct {
	ID            int                   `json:"id"`
	ShardKey      int                   `json:"shard_key"`
	Status        cluster.NodeStatus    `json:"status"`
	StoreNodeType cluster.StoreNodeType `json:"store_node_type"`
	LeaderID      int                   `json:"leader_id"`
	LeaderAddress struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"leader_address,omitempty"`
}

type ChangeLeaderRequest struct {
	NodeID int `json:"node_id" binding:"required"`
}

type ChangeLeaderResponse struct {
	Message   string `json:"message"`
	ShardID   int    `json:"shard_id"`
	OldLeader int    `json:"old_leader"`
	NewLeader int    `json:"new_leader"`
}
