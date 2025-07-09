package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	//r.Run(":8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// shutting down gracefully
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")

}
