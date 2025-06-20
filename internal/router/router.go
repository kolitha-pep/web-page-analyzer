package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kolitha-pep/web-page-analyzer/internal/handler"
)

func Setup() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/analyze/url", handler.AnalyzeHandler)
	}
	return r
}
