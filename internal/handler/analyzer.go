package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/analyzer"
	"github.com/sirupsen/logrus"
)

type urlAnalyzer struct {
	logger *logrus.Logger
}

type AnalyzerInterface interface {
	AnalyzeHandler(c *gin.Context)
}

func NewUrlAnalyzer(logger *logrus.Logger) AnalyzerInterface {
	return &urlAnalyzer{
		logger: logger,
	}
}

func (t urlAnalyzer) AnalyzeHandler(c *gin.Context) {
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

	t.logger.WithFields(logrus.Fields{
		"url":    url,
		"result": result,
	}).Info("Web page analysis completed")

	responseObject(c, result, nil)

}
