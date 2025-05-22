package service

import (
	"fmt"

	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/pkg/kvController/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type KvController struct {
	Router        *gin.Engine
	Config        *config.KvControllerConfig
	NodeManager   *NodeManager
	HealthManager *HealthManager
}

func NewKvController(cfg *config.KvControllerConfig) *KvController {
	logrus.Infof("Loaded controller config: %#v", cfg)

	controller := &KvController{}
	handler := api.NewRouteHandler(controller)
	router := api.SetupRouter(handler)

	controller.Router = router
	controller.Config = cfg
	return controller
}

func (c *KvController) Start() error {
	addr := c.Config.Address.Host + ":" + fmt.Sprint(c.Config.Address.Port)
	logrus.Infof("Starting KvController on %s", addr)
	return c.Router.Run(addr)
}

func (c *KvController) RegisterNode(address string, port int) error {
	return c.NodeManager.RegisterNode(address, port)
}

func (c *KvController) CheckNodesHealth() {
	c.HealthManager.checkNodes()
}

func (c *KvController) ChangePartitionLeader(shardID, targetNodeID int) error {
	shardInfo, exists := c.NodeManager.GetShardInfo(shardID)
	if !exists {
		return fmt.Errorf("shard %d not found", shardID)
	}

	var isFollower bool
	for _, f := range shardInfo.GetFollowers() {
		if f.GetID() == targetNodeID && f.GetStatus() == internal.NodeStatusActive {
			isFollower = true
			break
		}
	}
	if !isFollower {
		return fmt.Errorf("target node is not an active follower of this shard")
	}

	targetNode, err := c.NodeManager.GetNodeInfo(targetNodeID)
	if err != nil || targetNode.Status != internal.NodeStatusActive {
		return fmt.Errorf("invalid or inactive target node")
	}

	oldLeaderID := shardInfo.GetMaster().GetID()

	// Update master
	if err := c.NodeManager.UpdateShardMaster(shardID, targetNodeID); err != nil {
		return err
	}

	if err := c.HealthManager.notifyNewLeader(&targetNode); err != nil {
		return fmt.Errorf("failed to notify new leader: %v", err)
	}

	var followers []*NodeInfo
	for _, f := range shardInfo.GetFollowers() {
		if node, err := c.NodeManager.GetNodeInfo(f.GetID()); err == nil {
			followers = append(followers, &node)
		}
	}

	if err := c.HealthManager.notifyFollowers(followers, &targetNode); err != nil {
		logrus.WithError(err).Warn("Failed to notify some followers about leader change")
	}

	logrus.WithFields(logrus.Fields{
		"shard_id":   shardID,
		"old_leader": oldLeaderID,
		"new_leader": targetNodeID,
	}).Info("Shard leader changed successfully")

	return nil
}
