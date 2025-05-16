package api

type AddNodeRequest struct {
	NodeID  string `json:"node_id" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type RemoveNodeRequest struct {
	NodeID string `json:"node_id" binding:"required"`
}

type ChangeLeaderRequest struct {
	NewLeaderNodeID string `json:"new_leader_node_id" binding:"required"`
}

type MovePartitionRequest struct {
	TargetNodeID string `json:"target_node_id" binding:"required"`
}
