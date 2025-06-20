package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorRsp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type SuccessRsp struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

func ErrorResponse(data any, err error) ErrorRsp {
	return ErrorRsp{
		Success: false,
		Message: err.Error(),
		Data:    data,
	}
}

func responseObject(c *gin.Context, out any, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(out, err))
	} else {
		var data any
		if out == nil {
			data = nil
		} else {
			data = out
		}
		c.JSONP(http.StatusOK, SuccessRsp{
			Success: true,
			Data:    data,
		})
	}
}
