package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ControllerRouteHandler interface to decouple route definitions from implementation
type ControllerRouteHandler interface {
	HealthHandler(ctx *gin.Context)

	AddNodeHandler(ctx *gin.Context)
	RemoveNodeHandler(ctx *gin.Context)

	IncreasePartitionsHandler(ctx *gin.Context)
	DecreasePartitionsHandler(ctx *gin.Context)

	ChangePartitionLeaderHandler(ctx *gin.Context)
	MovePartitionHandler(ctx *gin.Context)

	NodeRegisterHandler(ctx *gin.Context)
	GetNodeInfoHandler(ctx *gin.Context)
	GetClusterHandler(ctx *gin.Context)
}

// SetupRouter initializes Gin router with routes bound to provided handlers
func SetupRouter(h ControllerRouteHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "KvController API is running")
	})
	router.GET("/health", h.HealthHandler)
	admin := router.Group("/admin")
	{
		// Node management
		admin.POST("/nodes", h.AddNodeHandler)
		admin.DELETE("/nodes/:id", h.RemoveNodeHandler)
		admin.GET("/nodes/:id", h.GetNodeInfoHandler)

		// Partition management
		admin.POST("/partitions/increase", h.IncreasePartitionsHandler)
		admin.POST("/partitions/decrease", h.DecreasePartitionsHandler)
		admin.POST("/partitions/:id/leader", h.ChangePartitionLeaderHandler)
		admin.POST("/partitions/:id/move", h.MovePartitionHandler)
		admin.GET("/cluster", h.GetClusterHandler)
	}

	internal := router.Group("/internal")
	{
		internal.POST("/nodes/register", h.NodeRegisterHandler)
	}
	log.Println("Controller router setup complete, new nodes can connect via /internal/nodes/register")

	return router
}
