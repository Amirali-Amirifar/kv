package api

type NodeRegisterHandlerRequest struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type NodeRegisterHandlerResponse struct {
	Error string `json:"error"`
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
