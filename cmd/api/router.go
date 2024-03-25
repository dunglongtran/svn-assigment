package api

import (
	"SVN-interview/cmd/api/handlers"
	"SVN-interview/internal/common"
	"github.com/gin-gonic/gin"
)

func SetupRouter(appCtx *common.AppContext) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		handlers.PingHandler(c, appCtx)
	})

	//TODO: Add more routes here

	return router
}
