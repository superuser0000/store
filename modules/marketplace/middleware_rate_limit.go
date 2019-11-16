package marketplace

import (
    "net/http"
    "sync"

    "github.com/gocraft/web"
    "golang.org/x/time/rate"
)

// IPRateLimiter .
type IPRateLimiter struct {
    ips map[string]*rate.Limiter
    mu  *sync.RWMutex
    r   rate.Limit
    b   int
}

// NewIPRateLimiter .
func NewUsernameRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    i := &IPRateLimiter{
        ips: make(map[string]*rate.Limiter),
        mu:  &sync.RWMutex{},
        r:   r,
        b:   b,
    }

    return i
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddUsername(username string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()

    limiter := rate.NewLimiter(i.r, i.b)

    i.ips[username] = limiter

    return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(username string) *rate.Limiter {
    i.mu.Lock()
    limiter, exists := i.ips[username]

    if !exists {
        i.mu.Unlock()
        return i.AddUsername(username)
    }

    i.mu.Unlock()

    return limiter
}

var (
    limiter      = NewUsernameRateLimiter(2, 20000)
    abuseLimiter = NewUsernameRateLimiter(10, 20000)
)

func (c *Context) RateLimitMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
    if c.ViewUser == nil {
        next(w, r)
        return
    }
    limiter := limiter.GetLimiter(c.ViewUser.Username)
    if !limiter.Allow() {
        http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
        return
    }

    next(w, r)
}

func (c *Context) AbuseRateLimitMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
    if c.ViewUser == nil {
        next(w, r)
        return
    }
    limiter := abuseLimiter.GetLimiter(c.ViewUser.Username)
    if !limiter.Allow() {
        c.ViewUser.User.Banned = true
        c.ViewUser.User.Save()
        return
    }

    next(w, r)
}
