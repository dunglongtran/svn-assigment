package handlers

import (
	"SVN-interview/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PingHandler(c *gin.Context, appCtx *common.AppContext) {
	// Ví dụ, sử dụng appCtx.DB để truy vấn database hoặc appCtx.Cache để tương tác với Redis

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
