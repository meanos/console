package app

import (
	"meanos.io/console/payment"
	"meanos.io/console/publisher"
	"meanos.io/console/review"
)

type Application struct {
	Name             string              `bson:"name"                 json:"name"`
	Description      string              `bson:"description"          json:"description"`
	IconURL          string              `bson:"icon_url"             json:"icon_url"`
	CoverURL         string              `bson:"cover_url"            json:"cover_url"`
	PrettyPreviewURL string              `bson:"pretty_preview_url"   json:"pretty_preview_url"`
	ScreenshotsURLs  []string            `bson:"screenshots_urls"     json:"screenshots_urls"`
	Rating           float64             `bson:"rating"               json:"rating"`
	PaymentType      payment.PaymentType `bson:"payment_type"         json:"payment_type"`

	Publisher publisher.Publisher `bson:"publisher"            json:"publisher"`

	Version   ApplicationVersion `bson:"version"              json:"version"`
	AppId     string             `bson:"app_id"               json:"app_id"`
	Slug      string             `bson:"slug"                 json:"slug"`
	Downloads int                `bson:"downloads"            json:"downloads"`
	Reviews   []review.Review    `bson:"reviews"              json:"reviews"`
	Tags      []string           `bson:"tags"                 json:"tags"`
	Category  string             `bson:"category"             json:"category"`
	Revenue   float64            `bson:"revenue"              json:"revenue"`

	Status string `bson:"status"               json:"status"`
}

type ApplicationOwners struct {
	Owners []AppLicationOwner `bson:"application_owners"   json:"application_owners"`
	AppId  string             `bson:"app_id"               json:"app_id"`
}

type AppLicationOwner struct {
	UID                 string   `bson:"uid"                  json:"uid"`
	PreferredCategories []string `bson:"preferred_categories" json:"preferred_categories"`
	AppsIDs             []string `bson:"apps_ids"             json:"apps_ids"`
}

type ApplicationVersion struct {
	AppId          string         `bson:"app_id"              json:"app_id"`
	CurrentVersion VersionRecord  `bson:"current_version"     json:"current_version"`
	History        VersionHistory `bson:"history"             json:"history"`
}

type VersionRecord struct {
	AppId        string `bson:"app_id"              json:"app_id"`
	Version      string `bson:"version"             json:"version"`
	MaintainerID string `bson:"maintainer_id"       json:"maintainer_id"`
	ReleaseNotes string `bson:"release_notes"       json:"release_notes"`
	ReleaseIndex int    `bson:"release_index"       json:"release_index"`
	PackageURL   string `bson:"package_url"         json:"package_url"`
}

type VersionHistory struct {
	Versions []VersionRecord `bson:"versions"            json:"versions"`
}

type ExportedApplication struct {
	Name              string                  `json:"name"`
	Version           string                  `json:"version"`
	Publisher         string                  `json:"publisher"`
	PublisherVerified bool                    `json:"publisher_verified"`
	Downloads         int                     `json:"downloads"`
	Rating            float64                 `json:"rating"`
	ReleaseNotes      string                  `json:"release_notes"`
	IconURL           string                  `json:"icon_url"`
	CoverURL          string                  `json:"cover_url"`
	PrettyPreviewURL  string                  `json:"pretty_preview_url"`
	ScreenshotsURLs   []string                `json:"screenshots_urls"`
	Category          string                  `json:"category"`
	Reviews           []review.ExportedReview `json:"reviews"`
	Free              bool                    `json:"free"`
	Monthly           bool                    `json:"monthly"`
	Yearly            bool                    `json:"yearly"`
	Once              bool                    `json:"once"`
	Price             float64                 `json:"price"`
	PackageURL        string                  `json:"package_url"`
	AppId             string                  `json:"app_id"`
}

type ExportedApplications struct {
	Apps []ExportedApplication `json:"apps"`
}

type ImportedAppUpdate struct {
	AppID   string `json:"app_id"`
	Version string `json:"version"`
}

type ImportedAppUpdates struct {
	TokenID string              `json:"token_id"`
	Records []ImportedAppUpdate `json:"records"`
}

type ExportedAppUpdate struct {
	AppID     string `json:"app_id"`
	Version   string `json:"version"`
	UpdateURl string `json:"update_url"`
	Paid      bool   `json:"paid"`
}

type ExportedAppUpdates struct {
	Updates []ExportedAppUpdate `json:"updates"`
}
