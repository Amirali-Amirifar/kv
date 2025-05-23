package api

import (
	"bytes"
	"github.com/Amirali-Amirifar/kv/internal/types/api"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Del(key string) error
	//UpdateNodeData() error
}

type HTTPServer struct {
	svc    Service
	router *gin.Engine
}

func NewHTTPServer(svc Service) *HTTPServer {
	router := gin.Default() // Includes logger and recovery middleware
	return &HTTPServer{
		svc:    svc,
		router: router,
	}
}

func (s *HTTPServer) Serve(port int) error {
	s.registerRoutes()
	log.Printf("Listening to connections on HTTP, Port: %d\n", port)

	return s.router.Run(":" + strconv.Itoa(port))
}

func (s *HTTPServer) registerRoutes() {
	s.router.Use(func(c *gin.Context) {
		// Read the body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.WithError(err).Warn("Failed to read request body")
			c.Next()
			return
		}

		// Restore the io.ReadCloser to the original state for next handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Log the body as a string (you can limit length if you want)
		log.WithFields(log.Fields{
			"Service": "HTTP Server",
			"Path":    c.Request.URL.Path,
			"Body":    string(bodyBytes),
		}).Info("Registering HTTP request")

		c.Next()
	})
	s.router.POST("/get", s.handleGet)
	s.router.POST("/set", s.handleSet)
	s.router.POST("/del", s.handleDel)
	s.router.POST("/health", s.handleHealth)
}

// handleGet processes GET requests
func (s *HTTPServer) handleGet(c *gin.Context) {
	var req api.GetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := s.svc.Get(req.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}

// handleSet processes SET requests
func (s *HTTPServer) handleSet(c *gin.Context) {
	var req api.SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.svc.Set(req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleDel processes DEL requests
func (s *HTTPServer) handleDel(c *gin.Context) {
	var req api.DelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.svc.Del(req.Key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.DelResponse{})
}

func (s *HTTPServer) UpdateNodeData(c *gin.Context) (string, error) {
	//var NodeData cluster.ShardInfo
	//s.UpdateNodeData()
	return "", nil
}
func (s *HTTPServer) handleHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
