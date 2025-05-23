// Package api is shared API types for node communications
// including client to load balancer, load balancer to master and master to replica

package api

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
