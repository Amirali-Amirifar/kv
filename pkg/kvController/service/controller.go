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

func (c *KvController) ChangePartitionLeader(shardID int, targetNodeID int) error {
	shardInfo, exists := c.NodeManager.GetShardInfo(shardID)
	if !exists {
		return fmt.Errorf("shard %d not found", shardID)
	}
	targetNode, err := c.NodeManager.GetNodeInfo(targetNodeID)
	if err != nil {
		return fmt.Errorf("target node not found: %v", err)
	}
	if targetNode.Status != internal.NodeStatusActive {
		return fmt.Errorf("target node is not active")
	}
	if err := c.NodeManager.UpdateShardMaster(shardID, targetNodeID); err != nil {
		return err
	}
	if err := c.HealthManager.notifyNewLeader(&targetNode); err != nil {
		return fmt.Errorf("failed to notify new leader: %v", err)
	}
	followers := make([]*NodeInfo, 0)
	for _, f := range shardInfo.GetFollowers() {
		if node, err := c.NodeManager.GetNodeInfo(f.GetID()); err == nil {
			followers = append(followers, &node)
		}
	}
	if err := c.HealthManager.notifyFollowers(followers, &targetNode); err != nil {
		logrus.WithError(err).Warn("Failed to notify some followers about leader change")
	}

	return nil
}

func (c *KvController) GetNodeManager() interface {
	GetShardInfo(shardID int) (interface {
		GetMaster() interface {
			GetID() int
			GetAddress() (string, int)
			GetStatus() internal.NodeStatus
		}
		GetFollowers() []interface {
			GetID() int
			GetAddress() (string, int)
			GetStatus() internal.NodeStatus
		}
	}, bool)
} {
	return c.NodeManager
}
