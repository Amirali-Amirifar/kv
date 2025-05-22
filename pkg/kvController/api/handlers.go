package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type KvControllerInterface interface {
	RegisterNode(address string, port int) error
	ChangePartitionLeader(shardID int, nodeID int) error
}

// KvRouteHandler implement ControllerRouteHandler
type KvRouteHandler struct {
	controller KvControllerInterface
}

func NewRouteHandler(controller KvControllerInterface) *KvRouteHandler {
	return &KvRouteHandler{
		controller: controller,
	}
}

// HealthHandler Returns status of all nodes.
func (k *KvRouteHandler) HealthHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// AddNodeHandler Adds a node to a designated partition
func (k *KvRouteHandler) AddNodeHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// RemoveNodeHandler Removes a node from a designated partition
func (k *KvRouteHandler) RemoveNodeHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// IncreasePartitionsHandler Adds a new partition
func (k *KvRouteHandler) IncreasePartitionsHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// DecreasePartitionsHandler Removes a partition.
func (k *KvRouteHandler) DecreasePartitionsHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// ChangePartitionLeaderHandler changes the leader from a partition
func (k *KvRouteHandler) ChangePartitionLeaderHandler(ctx *gin.Context) {
	shardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shard ID"})
		return
	}

	var req struct {
		NodeID int `json:"node_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := k.controller.ChangePartitionLeader(shardID, req.NodeID); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not a") {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "leader changed successfully",
		"shard_id":   shardID,
		"new_leader": req.NodeID,
	})
}

// MovePartitionHandler Moves a partition
func (k *KvRouteHandler) MovePartitionHandler(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (k *KvRouteHandler) NodeRegisterHandler(ctx *gin.Context) {
	req := &NodeRegisterHandlerRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logrus.WithError(err).Error("Failed to bind NodeRegisterHandlerRequest")
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := k.controller.RegisterNode(req.Ip, req.Port)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to register node %s:%d", req.Ip, req.Port)
		ctx.Status(http.StatusConflict)
		return
	}

	ctx.Status(http.StatusOK)
}
