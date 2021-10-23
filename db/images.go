package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/masci/flickr"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"
)

var FLICKR_KEY string
var FLICKR_SECRET string

type FlickrClient struct {
	Mutex  sync.Mutex
	Client *flickr.FlickrClient
}

var FClient FlickrClient

const (
	MB = 1 << 20
)

func UploadImage(f *multipart.FileHeader, c *gin.Context, path string) string {
	fmt.Println("We are uploading")
	if f.Size >= 5*MB {
		fmt.Println("OVERSIZE")
		return ""
	}
	if !(strings.Contains(f.Filename, ".png")) {
		fmt.Println(f.Filename, "NOT CONTAINS")
		return ""
	}
	pseudo := strconv.Itoa(int(time.Now().UnixNano())) + ".png"
	rname := "/pictures/" + path + "/" + pseudo
	err := c.SaveUploadedFile(f, rname)
	if err != nil {
		return ""
	}
	return pseudo
}
