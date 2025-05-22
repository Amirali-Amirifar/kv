package api

import "github.com/Amirali-Amirifar/kv/internal"

type GetRequest struct {
	Key string `json:"key"`
}

type GetResponse struct {
	Value string `json:"value"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetResponse struct{}

type DelRequest struct {
	Key string `json:"key"`
}

type DelResponse struct{}

type UpdateStateRequest struct {
	State    internal.StoreNodeType `json:"state"`
	LeaderID int                    `json:"leader_id"`
}
