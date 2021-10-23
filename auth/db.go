package auth

import (
	"context"
	"meanos.io/console/utils"
	"fmt"
	"github.com/meanOs/AMS"
	beatrix "github.com/meanOs/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

type DevToken struct {
	Id  string `json:"id"`
	UID string `json:"uid"`
}

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

	C.LoadCookiesManager()
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
	collection := Client.Client.Database("Users").Collection(collectionName)
	Client.Mutex.Unlock()
	return true, collection
}

func GetUserIdByEmailAndPassword(login, password string) (bool, string) {
	if t, c := NewDBCollection("accounts"); t {
		hash := utils.Makehash(password)

		filter := bson.M{"email": login, "password": hash}

		var res AMS.Account
		fmt.Println(c.Name(), c.Indexes())
		if err := c.FindOne(context.Background(), filter).Decode(&res); err == nil {
			log.Println(err)
			return res.UID != "", res.UID
		} else {
			log.Println(err)
			return false, ""
		}
	} else {
		return false, ""
	}
}

func (c *CookiesManager) LoadCookiesManager() {
	if t, collection := NewDBCollection("cookies"); t {
		filter := bson.M{"index": 0}

		var result ExportedManager

		if collection.FindOne(context.Background(), filter).Decode(&result) == nil {
			c.m = result.Map
			return
		} else {
			c.m = make(map[string]string)
			return
		}
	} else {
		c.m = make(map[string]string)
		return
	}
}

func (ex *ExportedManager) Dump() {
	if t, collection := NewDBCollection("cookies"); t {
		if _, err := collection.InsertOne(context.Background(), ex); err == nil {
			return
		} else {
			log.Println(err)
			go beatrix.SendError("ERROR: CANNOT INSERT TO DB", "EXPORTEDMANAGER.DUMP")
			return
		}
	} else {
		return
	}

}
