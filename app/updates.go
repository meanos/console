package app

import (
	"meanos.io/console/db"
	"meanos.io/console/publisher"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

func (a *Application) PublishUpdate(vr VersionRecord) bool {
	a.Version.CurrentVersion = vr
	a.Version.History.Versions = append(a.Version.History.Versions, vr)
	return a.UpdateDB()
}

func NewUpdate(uid, appId, vIndex, vDesc string, vUp *multipart.FileHeader, c *gin.Context) string {
	if t, application := GetAppByID(appId); t {
		if t, _ := publisher.GetPublisherByUID(uid); t {
			var vr VersionRecord
			vr.AppId = appId
			vr.Version = vIndex
			vr.MaintainerID = uid
			vr.ReleaseNotes = vDesc
			vr.ReleaseIndex = application.Version.CurrentVersion.ReleaseIndex + 1
			var prefix string
			if application.PaymentType.Free {
				prefix = "free"
			} else {
				prefix = "paid"
			}
			if vr.PackageURL = db.UploadFile(vUp, c, "/packages/"+prefix+"/"+application.Slug, ".xapp"); vr.PackageURL == "" {
				return "Error uploading file"
			} else {
				if application.PublishUpdate(vr) {
					return ""
				} else {
					return "Error updating application record"
				}
			}
		} else {
			return "Error verifying you are developer"
		}
	} else {
		return "Error getting application"
	}
}
