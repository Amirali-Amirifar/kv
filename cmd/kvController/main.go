package main

import (
	"fmt"
	"github.com/Amirali-Amirifar/kv/internal/config"
	"github.com/Amirali-Amirifar/kv/pkg/kvController/service"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.WithField("Package", "KvController")

func runKvController(configPath string) {
	log.Info("Starting KvController")

	var cfg config.KvControllerConfig
	config.LoadConfig(configPath, "", &cfg)

	controller := service.NewKvController(&cfg)
	err := controller.Start()

	if err != nil {
		log.Fatal(err)
		return
	}
}
func main() {
	var configPath string

	rootCmd := &cobra.Command{
		Use:   "kvController",
		Short: "kv node orchestration tool",
		Run: func(cmd *cobra.Command, args []string) {
			runKvController(configPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./config/controller_config.yaml", "Path to config file")
	//rootCmd.MarkPersistentFlagRequired("config") // <-- this line makes it mandatory

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
