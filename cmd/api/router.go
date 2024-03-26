package api

import (
	"SVN-interview/cmd/api/handlers"
	"SVN-interview/docs"
	"SVN-interview/internal/common"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(appCtx *common.AppContext) *gin.Engine {
	router := gin.Default()
	//docs.SwaggerInfo.BasePath = "/"
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		handlers.PingHandler(c, appCtx)
	})

	//TODO: Add more routes here

	return router
}
