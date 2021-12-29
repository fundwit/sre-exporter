package doc

import (
	_ "sre-exporter/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterDocsAPI
// - /api/index.html
// - /api/doc.json
func RegisterDocsAPI(r *gin.Engine, middleWares ...gin.HandlerFunc) {
	g := r.Group("/api/*any", middleWares...)
	g.GET("", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
