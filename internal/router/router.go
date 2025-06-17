package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/something", func(context *gin.Context) {
			fmt.Println("something")
		})
	}
	return r
}
