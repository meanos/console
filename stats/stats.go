package stats

import (
	"context"
	"meanos.io/console/publisher"
	"go.mongodb.org/mongo-driver/bson"
)

type CountryRec map[string]int

type CountryStats struct {
	Total   CountryRec            `bson:"total"`
	TotalDT map[string]CountryRec `bson:"total_dt"`
}

type CompanyStats struct {
	PublisherID    string  `bson:"pub_id"`
	TotalDownloads int     `bson:"total_downloads"`
	TotalRevenue   float64 `bson:"total_revenue"`
	TotalComments  int     `bson:"total_comments"`
	TotalRatings   float64 `bson:"total_ratings"`
	TotalApps      int     `bson:"total_apps"`

	DownloadsDT map[string]int     `bson:"downloads_dt"`
	RevenueDT   map[string]float64 `bson:"revenue_dt"`
	RatingsDT   map[string]float64 `bson:"ratings_dt"`
	CommentsDT  map[string]int     `bson:"comments_dt"`

	Country CountryStats `bson:"country_rec"`
}

func (c *CompanyStats) Load(uid string) {
	if t, p := publisher.GetPublisherByUID(uid); t {
		if t, cs := GetCompanyStatsByID(p.UID); t {
			c.TotalRatings = cs.TotalRatings
			c.CommentsDT = cs.CommentsDT
			c.TotalComments = cs.TotalComments
			c.TotalRevenue = cs.TotalRevenue
			c.TotalApps = cs.TotalApps
			c.DownloadsDT = cs.DownloadsDT
			c.RevenueDT = cs.RevenueDT
			c.RatingsDT = cs.RatingsDT
			c.CommentsDT = cs.CommentsDT
			c.Country = cs.Country
		}
	}
	return
}

type AppAdditionalStats struct {
	AppId        string       `bson:"app_id"`
	TotalRevenue float64      `bson:"total_revenue"`
	Countries    CountryStats `bson:"country_rec"`
}

func (aas *AppAdditionalStats) Load(appid, uid string) {
	if publisher.VerifyPublisherOwnsApp(appid, uid) {
		if t, coll := NewDBCollection("additional"); t {
			var res AppAdditionalStats
			filter := bson.M{"app_id": appid}
			if err := coll.FindOne(context.Background(), filter).Decode(&res); err == nil {
				aas.TotalRevenue = res.TotalRevenue
				aas.Countries = res.Countries
			}
		}
	}
	return
}
