package server

import (
	"meanos.io/console/app"
	"meanos.io/console/auth"
	"meanos.io/console/publisher"
	"meanos.io/console/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func handleRoot(c *gin.Context) {
	if tid, err := c.Cookie("devid"); err == nil {
		log.Println(tid)
		if t, uid := auth.AuthenticateCookie(tid); t {
			if x, _ := publisher.GetPublisherByUID(uid); x {
				c.HTML(http.StatusOK, "index.html", web.RenderIndex(uid))
				return
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/create")
				c.Abort()
				return
			}
		}
	}
	fmt.Println("REDIRECTED")
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	c.Abort()
	return
}

func handleNewAppGET(c *gin.Context) {
	if tid, err := c.Cookie("devid"); err == nil {
		log.Println(tid)
		if t, uid := auth.AuthenticateCookie(tid); t {
			c.HTML(http.StatusOK, "create_app.html", web.RenderNewApplication(uid, "", ""))
			return
		}
	}
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
	return
}

func handleNewAppPOST(c *gin.Context) {
	if tid, err := c.Cookie("devid"); err == nil {
		log.Println(tid)
		if t, uid := auth.AuthenticateCookie(tid); t {
			appname := c.PostForm("name")
			appdesc := c.PostForm("description")
			appicon, errI := c.FormFile("app_icon")
			appcover, errC := c.FormFile("app_cover")
			appversion := c.PostForm("app_version")
			appvdesc := c.PostForm("version_description")
			appvp, errP := c.FormFile("upload_package")
			appscreens := c.PostForm("screenshots")
			apppt := c.PostForm("pt")
			apppr := c.PostForm("price")

			if appname == "" || appdesc == "" || appversion == "" || appvdesc == "" {
				fmt.Println(appname, appdesc, appversion, appvdesc, appvp, appscreens, apppt, apppr)
				log.Println("SOMETHING IS MISSING")
				c.Redirect(http.StatusTemporaryRedirect, "/newApp")
				c.Abort()
				return
			}

			// App is free
			if apppr == "" {
				if errP != nil {
					log.Println("Upload is empty")
					c.HTML(http.StatusOK, "create_app.html", web.RenderNewApplication(uid, "Error uploading file, is it empty?", ""))
					c.Abort()
					return
				}
				resp := app.CreateFreeApp(appname, appdesc, appscreens, appversion, appvdesc, uid, appicon, appcover, appvp, c)
				if resp != "" {
					c.HTML(http.StatusOK, "create_app.html", web.RenderNewApplication(uid, resp, ""))
					c.Abort()
					return
				} else {
					c.Redirect(http.StatusFound, "/")
					c.Abort()
					return
				}
			} else { //App is paid
				log.Println("Create paid app")
			}

			if errI != nil {
				log.Println("------ERR ICON")
				log.Println(errI)
			}
			if errC != nil {
				log.Println("------ERR COVER")
				log.Println(errC)
			}
			fmt.Println(appname, appdesc, appversion, appvdesc, appvp, appscreens, apppt, apppr)
			return
		}
	}
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
	return
}

func handlePushUpdate(c *gin.Context) {
	if cid, err := c.Cookie("devid"); err == nil {
		if t, uid := auth.AuthenticateCookie(cid); t {
			appId := c.PostForm("appId")
			vIndex := c.PostForm("vindex")
			vDesc := c.PostForm("vdesc")
			vUp, errU := c.FormFile("vup")
			if errU != nil {
				log.Println(errU)
				c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, "Error processing file"))
				return
			}
			if publisher.VerifyPublisherOwnsApp(appId, uid) {
				resp := app.NewUpdate(uid, appId, vIndex, vDesc, vUp, c)
				log.Println(resp)
				if resp == "" {
					c.Redirect(http.StatusFound, "/")
					c.Abort()
					return
				} else {
					c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, resp))
					return
				}
			} else {
				c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, "We can't verify you own this app"))
				return
			}
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	}
}

func handleApp(c *gin.Context) {
	appId := c.DefaultQuery("appId", "")

	if appId == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
		return
	}
	if cid, err := c.Cookie("devid"); err == nil {
		if t, uid := auth.AuthenticateCookie(cid); t {
			c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, ""))
			return
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	c.Abort()
	return

}

func handleLogin(c *gin.Context) {
	login := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")

	if login == "" || password == "" {
		log.Println("Empty params")
		c.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}

	log.Println(login, password)

	if t, uid := auth.GetUserIdByEmailAndPassword(login, password); t {
		c.SetCookie("devid", auth.NewCookie(uid), int(time.Now().Add(12*30*time.Hour).Unix()), "/", websiteURL, false, false)
		log.Println("Set a cookie")

		// Now let's check if he is a verified publisher

		if t, _ := publisher.GetPublisherByUID(uid); t {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		} else {
			c.Redirect(http.StatusFound, "/create")
			c.Abort()
		}
		return
	} else {
		log.Println("Error getting auth info")
		c.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
}

func handleLogout(c *gin.Context) {
	if cid, err := c.Cookie("devid"); err == nil {
		auth.RemoveCookie(cid)
	}
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	c.Abort()
	return
}

func handleCreateGET(c *gin.Context) {
	if cid, err := c.Cookie("devid"); err == nil {
		if t, _ := auth.AuthenticateCookie(cid); t {
			c.HTML(http.StatusOK, "create_team.html", gin.H{})
			return
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	c.Abort()
	return
}

func handleCreatePOST(c *gin.Context) {

	tname := c.PostForm("tname")
	tmail := c.PostForm("tmail")
	turl := c.PostForm("turl")
	taddr := c.PostForm("taddr")

	if tname == "" || tmail == "" || turl == "" || taddr == "" {
		fmt.Println(tname, tmail, turl, taddr)
		fmt.Println(c.Request.PostForm)
		c.Redirect(http.StatusFound, "/create")
		c.Abort()
		return
	} else {
		if tid, err := c.Cookie("devid"); err == nil {
			log.Println(tid)

			if t, uid := auth.AuthenticateCookie(tid); t {
				publisher.Create(tname, turl, taddr, tmail, uid)
				c.Redirect(http.StatusFound, "/")
				c.Abort()
				return
			} else {
				c.Status(http.StatusBadRequest)
				return
			}
		} else {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
	}
}

func handleCompanyGET(c *gin.Context) {
	if cid, err := c.Cookie("devid"); err == nil {
		if t, uid := auth.AuthenticateCookie(cid); t {
			c.HTML(http.StatusOK, "company.html", web.RenderCompanyPage(uid, ""))
			return
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	c.Abort()
	return
}

func handleCompanyPOST(c *gin.Context) {
	if cid, err := c.Cookie("devid"); err == nil {
		if t, uid := auth.AuthenticateCookie(cid); t {
			cname := c.PostForm("cname")
			cmail := c.PostForm("cmail")
			caddr := c.PostForm("caddr")
			cweb := c.PostForm("cweb")
			cIcon, errI := c.FormFile("cico")
			cCover, errC := c.FormFile("ccover")

			withIcon, withCover := errI == nil, errC == nil
			fmt.Println(cname, cmail, caddr, cweb, withIcon, withCover, errI, errC)
			c.HTML(http.StatusOK, "company.html", web.RenderCompanyPage(uid, publisher.CreateInfoUpdate(cname, cmail, caddr, cweb, uid, withIcon, withCover, cIcon, cCover, c)))
		}
	}
}
