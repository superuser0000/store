package marketplace

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gocraft/web"
)

var (
	authorizedUrls = map[string]bool{
		"/":                    true,
		"/api/auth/login":      true,
		"/api/auth/register":   true,
		"/app":                 true,
		"/app/download":        true,
		"/auth/login":          true,
		"/auth/register":       true,
		"/auth/recover":        true,
		"/botcheck":            true,
		"/captcha":             true,
		"/favicon.ico":         true,
		"/item-image":          true,
		"/item-category-image": true,
		"/user-avatar":         true,
	}
)

// AuthMiddleware is responsible for authenicating user and setting appropriate values to a context
func (c *Context) AuthMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	var (
		userUuid     string
		session      = sessionStore.Load(r.Request)
		oneMinute, _ = time.ParseDuration("20m")
		now          = time.Now()
	)

	// API authenticaion is performed based on "teken" parameter in query string
	// user uuid is exracted from a database record associated with token
	if len(r.URL.Query()["token"]) > 0 {
		apiSession, err := FindAPISessionByToken(r.URL.Query()["token"][0])
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.APISession = apiSession
		if !c.APISession.IsTwoFactorSession || (c.APISession.IsTwoFactorSession && c.APISession.IsSecondFactorCompleted) {
			userUuid = apiSession.UserUuid
		}
		if apiSession.User.Banned {
			http.NotFound(w, r.Request)
			return
		}
		// Otherwise UUID is extracted from a session
	} else {
		userUuid, _ = session.GetString("UserUuid")
	}

	allowPassthrough := false
	if authorizedUrls[r.URL.Path] {
		allowPassthrough = true
	}

	for prefix, _ := range authorizedUrls {
		if prefix != "/" {
			if strings.HasPrefix(r.URL.Path, prefix) {
				allowPassthrough = true
			}
		}
	}

	if userUuid != "" {
		user, _ := FindUserByUuid(userUuid)
		if user == nil || user.Banned {
			http.NotFound(w, r.Request)
			return
		}

		if user.LastLoginDate.Add(oneMinute).Before(now) {
			user.LastLoginDate = now
			go user.Save()
		}

		vu := user.ViewUser(user.Language, false)
		c.ViewUser = &vu

	} else if !allowPassthrough {
		http.Redirect(w, r.Request, "/auth/login", 302)
		return
	} else {
		next(w, r)
		return
	}

	// If user has a store associated with them (staff of a store), set it here
	store, _ := FindStoreByUserUuid(c.ViewUser.Uuid)
	if store != nil && store.Uuid != "" {
		viewStore := store.ViewStore(c.ViewUser.Language)
		c.ViewUserStore = &viewStore
		c.EnableStoreStaffUI = true

		if store.LastLoginDate.Add(oneMinute).Before(now) {
			store.LastLoginDate = now
			database.Model(&Store{}).Where("uuid = ?", store.Uuid).Update("last_login_date", now)
		}
	}

	next(w, r)
}

func (c *Context) BotCheckMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/item-image") || strings.HasPrefix(r.URL.Path, "/user-avatar") || strings.HasPrefix(r.URL.Path, "/item-category-image") {
		next(w, r)
		return
	}

	session := sessionStore.Load(r.Request)
	checkUuid, _ := session.GetString("BotCheckUuid")

	if checkUuid != "" {
		_, err := redisClient.Get(checkUuid).Result()
		if err == nil {
			next(w, r)
		} else if !strings.HasPrefix(r.URL.Path, "/botcheck") && !strings.HasPrefix(r.URL.Path, "/captcha") {
			http.Redirect(w, r.Request, "/botcheck?redirect="+template.URLQueryEscaper(r.URL.String()), 302)
		} else {
			next(w, r)
		}
	} else if !strings.HasPrefix(r.URL.Path, "/botcheck") && !strings.HasPrefix(r.URL.Path, "/captcha") {
		http.Redirect(w, r.Request, "/botcheck?redirect="+template.URLQueryEscaper(r.URL.String()), 302)
	} else {
		next(w, r)
	}
}

func (c *Context) AdminMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if !c.ViewUser.IsAdmin {
		http.NotFound(w, r.Request)
		return
	}
	next(w, r)
}

func (c *Context) StaffMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser.IsStaff || c.ViewUser.IsAdmin {
		next(w, r)
		return
	}
	http.NotFound(w, r.Request)
}
