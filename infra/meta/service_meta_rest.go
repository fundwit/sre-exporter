package meta

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterMetaRestAPI(r *gin.Engine, middleWares ...gin.HandlerFunc) {
	g := r.Group("/", middleWares...)
	g.GET("", metaInfo)
}

// @ID get-meta-infomation
// @Success 200 {object} meta.ServiceInfo
// @Router / [get]
func metaInfo(c *gin.Context) {
	// localize.MustGetMessage("status-running")
	m := GetServiceMeta()
	c.JSON(http.StatusOK, &m)
}
