package server

import (
	"github.com/gin-gonic/gin"
)

var websiteURL string

func Init(r *gin.Engine, c Config) {
	websiteURL = c.WebsiteURL
	setApi(r)
	setPages(r)
}
