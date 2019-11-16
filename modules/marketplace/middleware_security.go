package marketplace

import (
	"net/http"

	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) SecurityHeadersMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	w.Header().Add("X-Frame-Options", "SAMEORIGIN")
	w.Header().Add("Content-Security-Policy", "script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; font-src 'self' data: ; default-src 'self'; frame-ancestors 'self'; disown-opener; form-action 'self'")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-Xss-Protection", "1; mode=block")
	w.Header().Add("Referrer-Policy", "no-referrer")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// w.Header().Add("access-control-allow-headers", "*")
	// w.Header().Add("access-control-allow-methods", "GET, POST, OPTIONS")

	next(w, r)
}

func (c *Context) OptionsHandler(w web.ResponseWriter, r *web.Request, methods []string) {
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Add("Access-Control-Max-Age", "100")
	w.Header().Add("access-control-allow-headers", "*")
}

func (c *Context) CSRFMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	// check csrf
	if r.Method == "POST" && c.ViewUser != nil {
		csrfToken, err := r.Cookie("csrf-token")
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}

		value, err := redisClient.Get(csrfToken.Value).Result()
		if err != nil || value != c.ViewUser.Uuid {
			http.NotFound(w, r.Request)
			return
		}

		redisClient.Del(csrfToken.Value)
	}

	if c.ViewUser != nil {
		csrfToken := util.GenerateUuid()
		redisClient.Set(csrfToken, c.ViewUser.Uuid, 0)

		cookie := http.Cookie{
			Name:   "csrf-token",
			Value:  csrfToken,
			MaxAge: 3600}

		http.SetCookie(w, &cookie)
	}

	next(w, r)
}

// func (c *Context) APICORSMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
// 	w.Header().Add("Access-Control-Allow-Origin", "*")
// 	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

// 	next(w, r)
// }

// var (
// 	rateLimiter = limiter.New(memory.NewStore(), limiter.Rate{
// 		Period: 1 * time.Hour,
// 		Limit:  1000,
// 	})
// )

// func (c *Context) ClientIDRateLimitMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

// 	cookie, _ := r.Cookie("client-uuid")

// 	if cookie == nil {
// 		cookie = http.Cookie{Name: "client-uuid", Value: util.GenerateUuid(), Path: "/", Expires: expire, MaxAge: 86400}

// 		http.SetCookie(w, cookie)
// 		next(w, r)
// 		return
// 	}

// 	rateLimiter.

// 	cookie.Value

// }
