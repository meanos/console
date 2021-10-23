package server

type Config struct {
	BeatrixToken   string `json:"beatrix-token"`
	BeatrixChannel string `json:"beatrix-channelID"`
	MongoURI       string `json:"mongo-uri"`
	CookieSecret   string `json:"cookie_secret"`
	WebsiteURL     string `json:"host_url"`
	FlickrApi      string `json:"flickr-id"`
	FlickrSecret   string `json:"flickr-secret"`
}
