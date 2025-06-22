package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type webView struct {
	logger *logrus.Logger
}

type WebViewInterface interface {
	AnalyzeUrlViewHandler(c *gin.Context)
}

func NewWebView(logger *logrus.Logger) WebViewInterface {
	return &webView{
		logger: logger,
	}
}

func (w *webView) AnalyzeUrlViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		// Render the HTML template for the web view
	})
}
