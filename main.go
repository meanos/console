package main

import (
	"encoding/json"
	"meanos.io/console/app"
	"meanos.io/console/auth"
	"meanos.io/console/db"
	"meanos.io/console/publisher"
	"meanos.io/console/server"
	"meanos.io/console/stats"
	"meanos.io/console/web"
	"fmt"
	"github.com/meanOs/AMS"
	beatrix "github.com/meanOs/Beatrix"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func LoadConfig() server.Config {
	var config server.Config
	configFile, err := os.Open("credentials.json")
	if err != nil {
		log.Panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	fmt.Println(config)
	err = configFile.Close()
	if err != nil {
		log.Panic(err)
	}
	return config
}

func Init(r *gin.Engine, c server.Config) {
	db.Init(c.MongoURI, c.FlickrApi, c.FlickrSecret)
	auth.Init(c.MongoURI)
	app.Init(c.MongoURI)
	stats.Init(c.MongoURI)
	publisher.Init(c.MongoURI)
	AMS.Init(c.MongoURI, "")
	beatrix.Init("CONSOLE", c.BeatrixToken, c.BeatrixChannel)
	server.Init(r, c)
	web.ServerName = "http://" + c.WebsiteURL
}

func main() {
	r := gin.Default()

	config := LoadConfig()

	Init(r, config)

	r.LoadHTMLGlob("static/*.html")
	r.Static("/assets", "./static/assets")

	// Those are needed paths for app icons and covers
	r.Static("/api/images/icons", "/pictures/icons")
	r.Static("/api/images/covers", "/pictures/covers")

	r.Run(":80")
}
