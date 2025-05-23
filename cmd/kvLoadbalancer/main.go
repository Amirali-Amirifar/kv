package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/pkg/kvLoadbalancer"
)

func main() {
	// Load configuration

	var cfg config.KvLoadBalancerConfig

	config.LoadConfig("config/loadbalancer_config.yaml", "", &cfg)

	// Initialize the loadbalancer service
	svc := kvLoadbalancer.NewLoadBalancerService(&cfg)

	svc.Serve()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
}
