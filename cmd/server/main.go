package main

import (
	"github.com/kolitha-pep/web-page-analyzer/internal/config"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/logger"
	"github.com/kolitha-pep/web-page-analyzer/internal/router"
)

func main() {

	// Initialize logger and load configuration
	config.LoadEnv()
	logger.Init()
	logger.Log.Info("Starting web page analyzer")

	r := router.Setup(logger.Log)
	r.Run(":8080")
}
