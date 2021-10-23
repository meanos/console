package stats

import (
	"context"
	"meanos.io/console/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
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
	collection := Client.Client.Database("stats").Collection(collectionName)
	Client.Mutex.Unlock()
	return true, collection
}

func GetCompanyStatsByID(publisherId string) (bool, CompanyStats) {
	if t, coll := NewDBCollection("company"); t {
		var res CompanyStats
		filter := bson.M{"pub_id": publisherId}
		if err := coll.FindOne(context.Background(), filter).Decode(&res); err == nil {
			return true, res
		} else {
			return false, CompanyStats{}
		}
	} else {
		return false, CompanyStats{}
	}
}

func (c *CompanyStats) GetCountriesDataSorted() (string, string, string, string, int, int, int, int, string, string, string, string) {
	if len(c.Country.Total) == 0 {
		return "", "", "", "", 0, 0, 0, 0, "", "", "", ""
	}

	keys := make([]string, 0, len(c.Country.Total))

	for k := range c.Country.Total {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var c1, c2, c3, c4, c1c, c2c, c3c, c4c string
	var c1d, c2d, c3d, c4d int

	for idx, k := range keys {
		if idx == 4 {
			break
		}

		var cn, cnc string
		var cnd int

		cn = k
		cnd = c.Country.Total[k]
		cnc = fmt.Sprintf("%f", (cnd-c.Country.TotalDT[k][utils.KeyOffset()])*100/cnd)

		if idx == 0 {
			c1, c1d, c1c = cn, cnd, cnc
			continue
		}
		if idx == 1 {
			c2, c2d, c2c = cn, cnd, cnc
			continue
		}
		if idx == 2 {
			c3, c3d, c3c = cn, cnd, cnc
			continue
		}
		if idx == 3 {
			c4, c4d, c4c = cn, cnd, cnc
			continue
		}
		continue
	}

	return c1, c2, c3, c4, c1d, c2d, c3d, c4d, c1c, c2c, c3c, c4c
}

func (c *CompanyStats) GetDDDataSorted() (int, int, int, int, int, int, int, int, int, int, int, int) {
	td := int(time.Now().Month())
	if td == 1 {
		return c.TotalDownloads, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	}

	var d1, d2, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12 int

	var i = 0
	for i <= 12 {
		var res int
		oft, m := utils.KeySetOffset(i)
		if oft == "" {
			res = 0
		}
		res = c.DownloadsDT[oft]

		switch m {
		case "December":
			d12 = res
		case "November":
			d11 = res
		case "October":
			d10 = res
		case "Septmber":
			d9 = res
		case "August":
			d8 = res
		case "July":
			d7 = res
		case "June":
			d6 = res
		case "May":
			d5 = res
		case "April":
			d4 = res
		case "March":
			d3 = res
		case "February":
			d2 = res
		case "January":
			d1 = res
		}
		i++
	}

	return d1, d2, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12
}

func (c *CompanyStats) GetCDDataSorted() (int, int, int, int, int, int, int, int, int, int, int, int) {
	td := int(time.Now().Month())
	if td == 1 {
		return c.TotalDownloads, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	}

	var d1, d2, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12 int

	var i = 0
	for i <= 12 {
		var res int
		oft, m := utils.KeySetOffset(i)
		if oft == "" {
			res = 0
		}
		res = c.CommentsDT[oft]

		switch m {
		case "December":
			d12 = res
		case "November":
			d11 = res
		case "October":
			d10 = res
		case "Septmber":
			d9 = res
		case "August":
			d8 = res
		case "July":
			d7 = res
		case "June":
			d6 = res
		case "May":
			d5 = res
		case "April":
			d4 = res
		case "March":
			d3 = res
		case "February":
			d2 = res
		case "January":
			d1 = res
		}
		i++
	}

	return d1, d2, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12
}

func (aas *AppAdditionalStats) GetTopCountry() (string, int) {
	if len(aas.Countries.Total) == 0 {
		return "", 0
	}

	keys := make([]string, 0, len(aas.Countries.Total))

	for k := range aas.Countries.Total {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys[0], aas.Countries.Total[keys[0]]
}
