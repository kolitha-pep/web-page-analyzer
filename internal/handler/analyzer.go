package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/analyzer"
)

func AnalyzeHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		responseObject(c, nil, errors.New("url is empty"))
		return
	}

	result, err := analyzer.AnalyzeWebPage(url)
	if err != nil {
		responseObject(c, nil, err)
		return
	}

	responseObject(c, result, nil)
}
