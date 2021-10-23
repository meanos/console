package app

import (
	"meanos.io/console/db"
	"meanos.io/console/payment"
	"meanos.io/console/publisher"
	"meanos.io/console/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const (
	MB = 1 << 20
)

func (a *Application) CreateUID() {
	a.AppId = utils.MakeAppHash(a.Name + a.Description + strconv.Itoa(int(time.Now().UnixNano())))
}
func (a *Application) InitVersions(version, uid, rnotes, paurl string) {
	var vr = VersionRecord{
		AppId:        a.AppId,
		Version:      version,
		MaintainerID: uid,
		ReleaseNotes: rnotes,
		ReleaseIndex: 0,
		PackageURL:   paurl,
	}
	var x []VersionRecord
	x = append(x, vr)
	a.Version.AppId = a.AppId
	a.Version.CurrentVersion = vr
	a.Version.History = VersionHistory{Versions: x}
	return
}

func (a *Application) MakeSlugFree() bool {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		a.Slug = strconv.Itoa(int(time.Now().UnixNano()))
		return os.MkdirAll(filepath.Join("/packages/free", a.Slug), os.ModePerm) == nil
	}
	if _, err := os.Stat("/packages/free" + reg.ReplaceAllString(a.Name, "")); err == nil {
		a.Slug = reg.ReplaceAllString(a.Name, "")
		return os.MkdirAll(filepath.Join("/packages/free", a.Slug), os.ModePerm) == nil
	} else {
		a.Slug = strconv.Itoa(int(time.Now().UnixNano()))
		return os.MkdirAll(filepath.Join("/packages/free", a.Slug), os.ModePerm) == nil
	}
}

func (a *Application) MakeSlugPaid() bool {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		a.Slug = strconv.Itoa(int(time.Now().UnixNano()))
		return os.MkdirAll(filepath.Join("/packages/paid", a.Slug), os.ModePerm) == nil
	}
	if _, err := os.Stat("/packages/free" + reg.ReplaceAllString(a.Name, "")); err == nil {
		a.Slug = reg.ReplaceAllString(a.Name, "")
		return os.MkdirAll(filepath.Join("/packages/paid", a.Slug), os.ModePerm) == nil
	} else {
		a.Slug = strconv.Itoa(int(time.Now().UnixNano()))
		return os.MkdirAll(filepath.Join("/packages/paid", a.Slug), os.ModePerm) == nil
	}
}

func CreateFreeApp(name, description, screenshots, appVersion, versionDescription, uid string, appIcon, appCover, packageFile *multipart.FileHeader, c *gin.Context) string {
	if t, pub := publisher.GetPublisherByUID(uid); t {
		var a = Application{
			Name:        name,
			Description: description,
			Rating:      0,
			PaymentType: payment.PaymentType{
				Price:           0,
				Monthly:         false,
				Yearly:          false,
				Once:            false,
				Free:            true,
				SubscriptionUID: "",
			},
			Downloads: 0,
			Status:    "review",
			Slug:      "",
		}
		a.CreateUID()
		if !a.MakeSlugFree() {
			return "Internal error"
		}
		if packageFile.Size >= 150*MB {
			fmt.Println("OVERSIZE")
			return "Your file exceeds 150 mb, please contact team for further instructions"
		}
		if package_url := db.UploadFile(packageFile, c, "/packages/free/"+a.Slug, ".xapp"); package_url != "" {
			a.InitVersions(appVersion, uid, versionDescription, package_url)
			a.IconURL = db.UploadImage(appIcon, c, "icons")
			a.CoverURL = db.UploadImage(appCover, c, "covers")
			pub.Update(a.AppId)
			a.Release()
			return ""
		} else {
			return "Internal error uploading package"
		}
	} else {
		return "NP"
	}

}
