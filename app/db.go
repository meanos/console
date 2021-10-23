package app

import (
	"context"
	beatrix "github.com/meanOs/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"sync"
	"time"
)

var Client DBConn
var URI string

func Init(mongouri string) {
	URI = mongouri
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	Client.Mutex.Lock()
	Client.Client = client
	Client.Mutex.Unlock()
}

type DBConn struct {
	Mutex  sync.Mutex
	Client *mongo.Client
}

func (c *DBConn) Reload() {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	Client.Mutex.Lock()
	Client.Client = client
	Client.Mutex.Unlock()
}

func NewDBCollection(collectionName string) (bool, *mongo.Collection) {
	Client.Mutex.Lock()
	collection := Client.Client.Database("AppStore").Collection(collectionName)
	Client.Mutex.Unlock()
	return true, collection
}

func GetPaidAppURL(id string) (int, string, string) {
	if t, c := NewDBCollection("apps"); t {
		filter := bson.M{"app_id": id}
		var result Application

		err := c.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			log.Println(err)
			log.Println("APPID", id)
			go beatrix.SendError("Error finding application by id", "GETPAIDAPPURL")
			return http.StatusInternalServerError, "", ""
		} else {
			slug := result.Slug
			filePath := result.Version.CurrentVersion.PackageURL

			if slug == "" || filePath == "" {
				go beatrix.SendError("ERROR: SLUG OR FILEPATH EMPTY", "GETPAIDAPPURL")
				return http.StatusInternalServerError, "", ""
			} else {
				return http.StatusOK, slug, filePath
			}
		}
	} else {
		go beatrix.SendError("ERROR CREATING NEW COLLECTION", "GETPAIDAPPURL")
		return http.StatusInternalServerError, "", ""
	}
}

func GetAppURL(id string) (bool, string) {
	if t, c := NewDBCollection("apps"); t {
		filter := bson.M{"app_id": id}
		var result Application

		err := c.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			log.Println(err)
			log.Println("APPID", id)
			go beatrix.SendError("Error finding application by id", "GETPAIDAPPURL")
			return false, ""
		} else {
			if !result.PaymentType.Free {
				//Sly ass trying to pass paid app url to get a free copy
				return false, ""
			} else {
				return true, result.Version.CurrentVersion.PackageURL
			}
		}
	} else {
		go beatrix.SendError("ERROR CREATING NEW COLLECTION", "GETPAIDAPPURL")
		return false, ""
	}
}

func GetAppByID(appid string) (bool, Application) {
	if t, c := NewDBCollection("apps"); t {
		filter := bson.M{"app_id": appid}
		var res Application
		err := c.FindOne(context.Background(), filter).Decode(&res)

		if err != nil {
			return false, Application{}
		}

		return true, res
	} else {
		return false, Application{}
	}
}

func (a *Application) Release() {
	if t, c := NewDBCollection("apps"); t {
		_, err := c.InsertOne(context.Background(), a)
		if err != nil {
			log.Println(err)
			go beatrix.SendError("Error pushing app to db", "APP.PUSH")
		}
	}
	return
}

func (a *Application) UpdateDB() bool {
	if t, c := NewDBCollection("apps"); t {
		filter := bson.M{"app_id": a.AppId}
		update := bson.M{"$set": bson.M{"version": a.Version}}
		r, err := c.UpdateOne(context.Background(), filter, update)
		log.Println(r, err)

		if err != nil {
			log.Println(err)
		}
		return err == nil
	}
	return false
}
