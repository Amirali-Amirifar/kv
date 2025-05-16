package api

import "github.com/gin-gonic/gin"

type KvRouteHandler struct {
}

// HealthHandler Returns status of all nodes.
func (k KvRouteHandler) HealthHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// AddNodeHandler Adds a node to a designated partition
func (k KvRouteHandler) AddNodeHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// RemoveNodeHandler Removes a node from a designated partition
func (k KvRouteHandler) RemoveNodeHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// IncreasePartitionsHandler Adds a new partition
func (k KvRouteHandler) IncreasePartitionsHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// DecreasePartitionsHandler Removes a partition.
func (k KvRouteHandler) DecreasePartitionsHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// ChangePartitionLeaderHandler changes the leader from a partition
func (k KvRouteHandler) ChangePartitionLeaderHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// MovePartitionHandler Moves a partition
func (k KvRouteHandler) MovePartitionHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
