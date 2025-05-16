package main

import (
	"github.com/Amirali-Amirifar/kv/internal"
	"github.com/Amirali-Amirifar/kv/internal/api"
	"github.com/Amirali-Amirifar/kv/pkg/kvNode"
	log "github.com/sirupsen/logrus"
	"net"
)

func main() {
	log.SetLevel(log.DebugLevel)
	_log := log.WithFields(log.Fields{"package": "KvNode"})

	config := &kvNode.Config{
		ID:     0,
		Status: internal.NodeStatusActive,
		Address: net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 8000,
		},
		StoreNodeType: internal.NodeTypeMaster,
	}

	_log.Info("Starting the KvNodeService")
	service := kvNode.NewKvNodeService(config)

	server := api.NewHTTPServer(service)
	err := server.Serve(config.Address.Port)
	if err != nil {
		panic(err)
	}
	_log.Info("Exiting")
}
