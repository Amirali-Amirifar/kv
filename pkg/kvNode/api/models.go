package api

import (
	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
)

type UpdateStateRequest struct {
	State    cluster.StoreNodeType `json:"state"`
	LeaderID int                   `json:"leader_id"`
}

type WALProgressRequest struct {
	FollowerID int   `json:"follower_id"`
	Seq        int64 `json:"seq"`
}
