package service

import (
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/pkg/kvController/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type KvController struct {
	Router      *gin.Engine
	Config      *config.KvControllerConfig
	NodeManager *NodeManager
}

func NewKvController(configPath string) *KvController {
	cfg := &config.KvControllerConfig{}
	config.LoadConfig(configPath, "KvControllerConfig", cfg)

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
