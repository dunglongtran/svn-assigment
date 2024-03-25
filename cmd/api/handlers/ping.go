package handlers

import (
	"SVN-interview/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PingHandler(c *gin.Context, appCtx *common.AppContext) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
