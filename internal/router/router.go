package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kolitha-pep/web-page-analyzer/internal/handler"
	"github.com/sirupsen/logrus"
)

func Setup(logger *logrus.Logger) *gin.Engine {

	analyzeHandler := handler.NewUrlAnalyzer(logger)
	webViewHandler := handler.NewWebView(logger)

	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")

	// API routes starts here
	api := r.Group("/api")
	{
		analyzeApi := api.Group("/analyze")
		{
			analyzeApi.GET("/url", analyzeHandler.AnalyzeHandler)
		}
	}

	// Web routes starts here
	web := r.Group("/web")
	{
		web.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
		web.GET("/", webViewHandler.AnalyzeUrlViewHandler)
	}

	return r
}
