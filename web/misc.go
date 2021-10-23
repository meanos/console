package web

const (
	placeholderIcon  = "https://via.placeholder.com/110x110"
	placeholderCover = "https://via.placeholder.com/800x500"
)

func renderIcon(imgUrl string) string {
	if imgUrl == "" {
		return placeholderIcon
	} else {
		return ServerName + "/api/images/icons/" + imgUrl
	}
}

func renderCover(imgUrl string) string {
	if imgUrl == "" {
		return placeholderCover
	} else {
		return ServerName + "/api/images/covers/" + imgUrl
	}
}
