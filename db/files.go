package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

var accessKeyID, secretAccessKey string

func uploadToSpace(fname, fpath string) (bool, string) {
	endpoint := "sfo2.digitaloceanspaces.com"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := "mean-updates"

	contentType := "application/zip"

	// Upload the zip file with FPutObject
	x, err := minioClient.FPutObject(bucketName, fname, fpath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println(err)
		return false, ""
	}
	log.Println(x)
	return true, fname
}

func UploadApplication(m *multipart.FileHeader, c *gin.Context) string {
	fmt.Println("We are uploading")

	if !(strings.Contains(m.Filename, ".xapp")) {
		fmt.Println(m.Filename, "NOT CONTAINS")
		return ""
	}
	pseudo := strconv.Itoa(int(time.Now().UnixNano())) + ".xapp"
	rname := "/temp/apps/" + pseudo
	err := c.SaveUploadedFile(m, rname)
	if err != nil {
		return ""
	}

	return pseudo
}

func UploadFile(m *multipart.FileHeader, c *gin.Context, path string, ftype string) string {
	fmt.Println("We are uploading")

	if !(strings.Contains(m.Filename, ftype)) {
		fmt.Println(m.Filename, "NOT CONTAINS")
		return ""
	}
	pseudo := strconv.Itoa(int(time.Now().UnixNano())) + ftype
	rname := path + "/" + pseudo
	err := c.SaveUploadedFile(m, rname)
	if err != nil {
		return ""
	}
	return pseudo
}
