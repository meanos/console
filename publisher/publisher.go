package publisher

import (
	"context"
	"meanos.io/console/db"
	"meanos.io/console/utils"
	beatrix "github.com/meanOs/Beatrix"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"mime/multipart"
	"strconv"
	"time"
)

type Publisher struct {
	DisplayName     string   `bson:"display_name"     json:"display_name"`
	MaintainersUIDs []string `bson:"maintainers_uids" json:"maintainers_uids"`
	SumRatings      float64  `bson:"sum_ratings"      json:"sum_ratings"`
	Email           string   `bson:"email"            json:"email"`
	Address         string   `bson:"address"          json:"address"`
	Website         string   `bson:"website"          json:"website"`
	Verified        bool     `bson:"verified"         json:"verified"`
	UID             string   `bson:"uid"              json:"uid"`
	Apps            []string `bson:"apps"             json:"apps"`

	ProfileIcon string `bson:"profile_icon"     json:"profile_icon" `
	CoverImage  string `bson:"cover_image"      json:"cover_image" `
}

func Create(tname, turl, taddr, tmail, uid string) {
	var p Publisher

	p.DisplayName = tname
	p.Website = turl
	p.Email = tmail
	p.Address = taddr
	p.MaintainersUIDs = append(p.MaintainersUIDs, uid)

	p.Verified = false
	p.UID = utils.Makehash(tname + turl + taddr + tmail + uid + strconv.Itoa(int(time.Now().UnixNano())))

	if t, c := NewDBCollection("publishers"); t {
		_, err := c.InsertOne(context.Background(), p)

		if err != nil {
			log.Println()
			go beatrix.SendError("Error inserting new collection", "PUBLISHER/CREATE")
		}
		return
	}
	return
}

func (p *Publisher) Update(appId string) {
	p.Apps = append(p.Apps, appId)
	filter := bson.M{"uid": p.UID}
	update := bson.M{"$set": bson.M{"apps": p.Apps}}
	if t, c := NewDBCollection("publishers"); t {
		if _, err := c.UpdateOne(context.Background(), filter, update); err == nil {
			return
		} else {
			log.Println(err)
			return
		}
	}
	return
}

func (p *Publisher) PushInfoUpdate() string {
	filter := bson.M{"uid": p.UID}
	update := bson.M{"$set": p}
	if t, c := NewDBCollection("publishers"); t {
		if _, err := c.UpdateOne(context.Background(), filter, update); err == nil {
			return ""
		} else {
			log.Println(err)
			return "Internal error"
		}
	}
	return "Internal error"
}

func CreateInfoUpdate(name, email, address, website, uid string, withIcon, withCover bool, icon, cover *multipart.FileHeader, c *gin.Context) string {
	if t, p := GetPublisherByUID(uid); t {
		if name != "" {
			p.DisplayName = name
		}
		if email != "" {
			p.DisplayName = name
		}
		if address != "" {
			p.DisplayName = name
		}
		if website != "" {
			p.DisplayName = name
		}
		if withIcon {
			if resp := db.UploadFile(icon, c, "/pictures/icons", ".png"); resp != "" {
				p.ProfileIcon = resp
			}
		}
		if withCover {
			if resp := db.UploadFile(cover, c, "/pictures/covers", ".png"); resp != "" {
				p.CoverImage = resp
			}
		}
		return p.PushInfoUpdate()
	} else {
		return "NP"
	}
}
