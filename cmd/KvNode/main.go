package main

import (
	"fmt"
	"os"

	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/pkg/kvNode"
	"github.com/Amirali-Amirifar/kv/pkg/kvNode/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runKvNode(configPath string) {
	log.Info("Starting kvNode")
	var cfg config.KvNodeConfig
	config.LoadConfig(configPath, "", &cfg)
	log.Printf("Loaded Config %#v", cfg)
	service := kvNode.NewKvNodeService(&cfg)

	server := api.NewHTTPServer(service)
	err := server.Serve(cfg.Address.Port)
	if err != nil {
		panic(err)
	}
	log.Info("Exiting")
}
func main() {

	var configPath string

	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "KV Address Service",
		Run: func(cmd *cobra.Command, args []string) {
			runKvNode(configPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./node_config.yaml", "Path to config file")
	//rootCmd.MarkPersistentFlagRequired("config") // <-- this line makes it mandatory

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
