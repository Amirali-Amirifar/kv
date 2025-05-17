package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type KvService interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Del(key string) error
}

type HTTPServer struct {
	svc    KvService
	router *gin.Engine
}

// NewHTTPServer creates a new HTTP server with the given KV service
func NewHTTPServer(svc KvService) *HTTPServer {
	router := gin.Default() // Includes logger and recovery middleware
	return &HTTPServer{
		svc:    svc,
		router: router,
	}
}

// Serve starts the HTTP server on the specified port
func (s *HTTPServer) Serve(port int) error {
	s.registerRoutes()
	log.Printf("Listening to connections on HTTP, Port: %d\n", port)

	return s.router.Run(":" + strconv.Itoa(port))
}

// registerRoutes sets up the API endpoints
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
			"Body":    fmt.Sprintf(string(bodyBytes)),
		}).Info("Registering HTTP request")

		c.Next()
	})
	s.router.POST("/get", s.handleGet)
	s.router.POST("/set", s.handleSet)
	s.router.POST("/del", s.handleDel)
}

// handleGet processes GET requests
func (s *HTTPServer) handleGet(c *gin.Context) {
	var req GetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, err := s.svc.Get(req.Key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetResponse{Value: val})
}

// handleSet processes SET requests
func (s *HTTPServer) handleSet(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.svc.Set(req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SetResponse{})
}

// handleDel processes DEL requests
func (s *HTTPServer) handleDel(c *gin.Context) {
	var req DelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.svc.Del(req.Key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, DelResponse{})
}
