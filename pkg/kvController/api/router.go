package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
}

// SetupRouter initializes Gin router with routes bound to provided handlers
func SetupRouter(h ControllerRouteHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", h.HealthHandler)

	router.POST("/nodes", h.AddNodeHandler)
	router.DELETE("/nodes/:id", h.RemoveNodeHandler)

	router.POST("/partitions/increase", h.IncreasePartitionsHandler)
	router.POST("/partitions/decrease", h.DecreasePartitionsHandler)

	router.POST("/partitions/:id/leader", h.ChangePartitionLeaderHandler)
	router.POST("/partitions/:id/move", h.MovePartitionHandler)

	// Optional: simple root endpoint
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Controller API is running")
	})

	return router
}
