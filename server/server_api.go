package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setApi(r *gin.Engine) {
	r.GET("/alive", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.GET("/apps", handleApiApps)
}
