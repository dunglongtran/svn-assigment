package handlers

import (
	"SVN-interview/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /

// PingHandler godoc
// @Summary Show the status of the API
// @Description do ping
// @Tags healthcheck
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func PingHandler(c *gin.Context, appCtx *common.AppContext) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
