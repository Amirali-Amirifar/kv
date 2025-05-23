package api

import (
	"github.com/Amirali-Amirifar/kv/internal/types"
)

type UpdateStateRequest struct {
	State    types.StoreNodeType `json:"state"`
	LeaderID int                 `json:"leader_id"`
}
