package auth

import (
	"meanos.io/console/utils"
	"strconv"
	"time"
)

var C CookiesManager

func AuthenticateCookie(cookieID string) (bool, string) {
	return C.Load(cookieID)
}

func AuthenticatePassword(login, password string) (bool, string) {
	return GetUserIdByEmailAndPassword(login, password)
}

func NewCookie(uid string) string {
	t := utils.Makehash(uid + strconv.Itoa(int(time.Now().Unix())))
	C.Write(t, uid)
	return t
}

func RemoveCookie(cid string) {
	C.Remove(cid)
}
