package main

import (
	"log"

	"github.com/kolitha-pep/web-page-analyzer/internal/router"
)

func main() {
	//config.LoadEnv()
	//db.Connect()
	r := router.Setup()
	log.Fatal(r.Run(":8080"))
}
