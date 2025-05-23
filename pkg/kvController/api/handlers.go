package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Amirali-Amirifar/kv/internal/types/cluster"
	"github.com/Amirali-Amirifar/kv/pkg/kvController/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// KvRouteHandler implement ControllerRouteHandler
type KvRouteHandler struct {
	controller interfaces.KvControllerInterface
}

func NewRouteHandler(controller interfaces.KvControllerInterface) *KvRouteHandler {
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

	var req ChangeLeaderRequest
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

	nodeInfo, err := k.controller.RegisterNode(req.Ip, req.Port)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to register node %s:%d", req.Ip, req.Port)
		ctx.Status(http.StatusConflict)
		return
	}

	response := NodeRegisterHandlerResponse{
		ID:            nodeInfo.ID,
		ShardKey:      nodeInfo.ShardKey,
		Status:        nodeInfo.Status,
		StoreNodeType: nodeInfo.StoreNodeType,
		LeaderID:      nodeInfo.LeaderID,
	}

	// If this is a follower, get the master's address
	if nodeInfo.StoreNodeType == cluster.NodeTypeFollower {
		masterInfo, err := k.controller.GetNodeManager().GetNodeInfo(nodeInfo.LeaderID)
		if err != nil {
			logrus.WithError(err).Error("Failed to get master node info")
			ctx.Status(http.StatusInternalServerError)
			return
		}
		response.LeaderAddress.IP = masterInfo.Address.IP.String()
		response.LeaderAddress.Port = masterInfo.Address.Port
	}

	ctx.JSON(http.StatusOK, response)
}

// GetNodeInfoHandler returns information about a specific node
func (k *KvRouteHandler) GetNodeInfoHandler(ctx *gin.Context) {
	nodeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
		return
	}

	nodeInfo, err := k.controller.GetNodeManager().GetNodeInfo(nodeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": nodeInfo.ID,
		"address": gin.H{
			"ip":   nodeInfo.Address.IP.String(),
			"port": nodeInfo.Address.Port,
		},
	})
}

func (k *KvRouteHandler) GetClusterHandler(ctx *gin.Context) {
	nodes := k.controller.GetClusterDetails()
	shardMap := make(map[int][]gin.H)
	for _, node := range nodes {
		nodeInfo := gin.H{
			"id":        node.ID,
			"shard_key": node.ShardKey,
			"status":    node.Status,
			"node_type": node.StoreNodeType,
			"leader_id": node.LeaderID,
			"address": gin.H{
				"ip":   node.Address.IP.String(),
				"port": node.Address.Port,
			},
		}
		shardMap[node.ShardKey] = append(shardMap[node.ShardKey], nodeInfo)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"shards": shardMap,
	})
}
