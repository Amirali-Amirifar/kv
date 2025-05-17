package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type KvControllerInterface interface {
	RegisterNode(address string, port int) error
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
	//TODO implement me
	panic("implement me")
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
