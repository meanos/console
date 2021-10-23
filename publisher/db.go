package publisher

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	collection := Client.Client.Database("dev").Collection(collectionName)
	Client.Mutex.Unlock()
	return true, collection
}

//User's id
func GetPublisherByUID(uid string) (bool, Publisher) {
	if t, c := NewDBCollection("publishers"); t {
		filter := bson.M{"maintainers_uids": uid}

		fmt.Println(filter)

		var res Publisher

		return c.FindOne(context.Background(), filter).Decode(&res) == nil, res
	} else {
		return false, Publisher{}
	}
}

func VerifyPublisherOwnsApp(appid, userid string) bool {
	if t, p := GetPublisherByUID(userid); t {
		for _, a := range p.Apps {
			if a == appid {
				return true
			}
		}
	}
	log.Println(appid, userid)
	return false
}

func GetAppIds(uid string) (bool, []string) {
	if t, c := NewDBCollection("publishers"); t {
		filter := bson.M{
			"maintainers_uids": uid,
		}
		var res Publisher
		return c.FindOne(context.Background(), filter).Decode(&res) == nil, res.Apps
	}
	return false, []string{}
}
