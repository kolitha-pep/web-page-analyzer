package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kolitha-pep/web-page-analyzer/internal/handler"
	"github.com/sirupsen/logrus"
)

func Setup(logger *logrus.Logger) *gin.Engine {

	analyzeHandler := handler.NewUrlAnalyzer(logger)

	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/analyze/url", analyzeHandler.AnalyzeHandler)
	}
	return r
}
