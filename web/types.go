package web

import (
	"meanos.io/console/app"
	"meanos.io/console/publisher"
	"meanos.io/console/stats"
	"meanos.io/console/utils"
	"fmt"
	"github.com/meanOs/AMS"
	"log"
	"strconv"
	"sync"
)

// Company stats prepared
type CSPrepared struct {
	TD  string
	TDC string
	TDW string
	AD  string

	TR  string
	TRC string
	TRW string
	AR  string

	TRT  string
	TRTC string
	TRTW string
	ART  string

	TC  string
	TCC string
	TCW string
	AC  string

	DDJA int
	DDFE int
	DDMA int
	DDAP int
	DDMY int
	DDJN int
	DDJL int
	DDAU int
	DDSE int
	DDOC int
	DDNO int
	DDDE int

	DCJA int
	DCFE int
	DCMA int
	DCAP int
	DCMY int
	DCJN int
	DCJL int
	DCAU int
	DCSE int
	DCOC int
	DCNO int
	DCDE int

	C1 string
	C2 string
	C3 string
	C4 string

	C1D int
	C2D int
	C3D int
	C4D int

	C1C string
	C2C string
	C3C string
	C4C string
}

func (c *CSPrepared) Load(cid string, wg *sync.WaitGroup) {
	var cs stats.CompanyStats
	cs.Load(cid)

	if cs.TotalApps == 0 {
		c.AD, c.AR, c.ART, c.AC = "up", "up", "up", "up"
		c.C1, c.C2, c.C3, c.C4 = "No data", "No data", "No data", "No data"
		c.C1C, c.C2C, c.C3C, c.C4C = "No data", "No data", "No data", "No data"
		c.TD, c.TR, c.TRT, c.TC = "No data", "No data", "No data", "No data"
		wg.Done()
		return
	}

	c.TD = strconv.Itoa(cs.TotalDownloads)
	tdc := (cs.TotalDownloads - cs.DownloadsDT[utils.KeyOffset()]) * 100 / (cs.DownloadsDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TDC = fmt.Sprintf("%f", tdc)
	if tdc > 0 {
		c.TDW = c.TDC
		c.AD = "up"
	} else {
		c.TDW = "0"
		c.AD = "down"
	}

	c.TR = fmt.Sprintf("%f", cs.TotalRevenue)
	var trc = ((cs.TotalRevenue - cs.RevenueDT[utils.KeyOffset()]) * 100) / (cs.RevenueDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TRC = fmt.Sprintf("%f", trc)
	if trc > 0 {
		c.TRW = strconv.Itoa(int(trc))
		c.AD = "up"
	} else {
		c.TRW = "0"
		c.AD = "down"
	}

	c.TRT = fmt.Sprintf("%f", cs.TotalRatings)
	TRTc := (cs.TotalRatings - cs.RatingsDT[utils.KeyOffset()]) * 100 / (cs.RatingsDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TRTC = fmt.Sprintf("%f", TRTc)
	if TRTc > 0 {
		c.TRTW = strconv.Itoa(int((cs.TotalRatings / 5) * 100))
		c.AD = "up"
	} else {
		c.TRTW = "0"
		c.AD = "down"
	}

	c.TC = strconv.Itoa(cs.TotalComments)
	TCc := (cs.TotalComments - cs.CommentsDT[utils.KeyOffset()]) * 100 / (cs.CommentsDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TCC = fmt.Sprintf("%f", TCc)
	if TCc > 0 {
		c.TCW = c.TCC
		c.AD = "up"
	} else {
		c.TCW = "0"
		c.AD = "down"
	}

	c.C1, c.C2, c.C3, c.C4, c.C1D, c.C2D, c.C3D, c.C4D, c.C1C, c.C2C, c.C3C, c.C4C = cs.GetCountriesDataSorted()

	c.DDJA, c.DDFE, c.DDMA, c.DDAP, c.DDMY, c.DDJN, c.DDJL, c.DDAU, c.DDSE, c.DDOC, c.DDNO, c.DDDE = cs.GetDDDataSorted()

	c.DCJA, c.DCFE, c.DCMA, c.DCAP, c.DCMY, c.DCJN, c.DCJL, c.DCAU, c.DCSE, c.DCOC, c.DCNO, c.DCDE = cs.GetCDDataSorted()

	wg.Done()
	return
}

type Account struct {
	Name   string
	Email  string
	PicURL string
	UID    string
}

func (a *Account) Load(uid string, wg *sync.WaitGroup) {
	if acc := AMS.GetUserByID(uid); acc.UID != "" {
		a.Name = acc.Name
		a.Email = acc.Email
		a.PicURL = acc.AvatarURL
	}
	wg.Done()
	return
}

type AppStatsExported struct {
	AppName         string
	AppDescription  string
	AppRating       float64
	AppRatingWidth  string
	AppDownloads    int
	AppRevenue      string
	AppComments     int
	ExportedCountry string
	UpdateType      string
	ID              string
}

func (ase *AppStatsExported) Load(appId, uid string, wg *sync.WaitGroup) {
	if t, application := app.GetAppByID(appId); t {
		var aas stats.AppAdditionalStats
		aas.Load(appId, uid)
		ase.AppName = application.Name
		ase.AppDescription = application.Description
		ase.AppRating = application.Rating
		ase.AppRatingWidth = strconv.Itoa(int((ase.AppRating / 5) * 100))
		ase.AppRevenue = "$" + fmt.Sprintf("%f", aas.TotalRevenue)
		ase.AppDownloads = application.Downloads
		ase.AppComments = len(application.Reviews)
		tc, tcd := aas.GetTopCountry()
		ase.ExportedCountry = tc + " (" + strconv.Itoa(int((tcd/(ase.AppDownloads+1))*100)) + "%)"
		if !application.PaymentType.Free {
			ase.UpdateType = "file"
		} else {
			ase.UpdateType = "url"
		}
		ase.ID = application.AppId
	}
	wg.Done()
	return
}

type AppTableElement struct {
	AppName        string `json:"app_name"`
	IconUrl        string `json:"icon_url"`
	DownloadsCount string `json:"downloads_count"`
	Revenue        string `json:"revenue"`
	LastUpdate     string `json:"last_update"`
	RatingWidth    string `json:"rating_width"`
	AppId          string `json:"app_id"`
}

type AppTable struct {
	Apps []AppTableElement `json:"apps"`
}

func (apt *AppTable) Load(uid string) {
	if t, apps := publisher.GetAppIds(uid); t {
		log.Println(apps)
		if len(apps) != 0 {
			for _, a := range apps {
				if t, apx := app.GetAppByID(a); t {
					var ate AppTableElement
					review := ""
					if apx.Status != "" {
						review = " (on review)"
					}
					ate.AppName = apx.Name + review
					ate.IconUrl = apx.IconURL
					ate.DownloadsCount = strconv.Itoa(apx.Downloads)
					ate.LastUpdate = apx.Version.CurrentVersion.Version
					ate.RatingWidth = fmt.Sprintf("%f", (apx.Rating/5)*100)
					ate.AppId = a
					ate.Revenue = fmt.Sprintf("%f", apx.Revenue)
					apt.Apps = append(apt.Apps, ate)
				} else {
					log.Println("Skipped")
				}
			}
		} else {
			log.Println("Is null")
		}
	} else {
		log.Println("T is false")
	}

	return
}

//Exported publisher
type EP struct {
	Name         string
	Email        string
	Address      string
	Website      string
	CompanyIcon  string
	CompanyCover string
}

func (ep *EP) Load(uid string, wg *sync.WaitGroup) {
	if t, pub := publisher.GetPublisherByUID(uid); t {
		ep.Name = pub.DisplayName
		ep.Website = pub.Website
		ep.Email = pub.Email
		ep.Address = pub.Address
		ep.CompanyIcon = pub.ProfileIcon
		ep.CompanyCover = pub.CoverImage
	}
	wg.Done()
	return
}
