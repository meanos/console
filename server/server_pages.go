package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setPages(r *gin.Engine) {
	r.GET("/", handleRoot)

	r.GET("/newApp", handleNewAppGET)

	r.POST("/newApp", handleNewAppPOST)

	r.POST("/pushUpdate", handlePushUpdate)

	r.GET("/app", handleApp)

	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", "") })

	r.POST("/login", handleLogin)

	r.GET("/logout", handleLogout)

	r.GET("/create", handleCreateGET)

	r.POST("/create", handleCreatePOST)

	r.GET("/company", handleCompanyGET)

	r.POST("/company", handleCompanyPOST)
}
