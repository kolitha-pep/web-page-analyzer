package main

import (
	"github.com/kolitha-pep/web-page-analyzer/internal/config"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/logger"
	"github.com/kolitha-pep/web-page-analyzer/internal/router"
)

func init() {
	config.LoadEnv()
}
func main() {

	// Initialize logger
	logger.Init()
	logger.Log.Info("Starting web page analyzer")

	r := router.Setup(logger.Log)
	r.Run(":8080")
}
