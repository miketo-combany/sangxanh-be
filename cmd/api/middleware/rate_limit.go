package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type ipEntry struct {
	count     int
	lastVisit time.Time
}

var (
	rateLimitMap = make(map[string]*ipEntry)
	mu           sync.Mutex
)

// Configurable constants
const (
	maxRequestsPerWindow = 10
	rateWindow           = 1 * time.Minute
	blockDuration        = 10 * time.Minute
)

var blockedIPs = make(map[string]time.Time)

func RateLimiterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()

		mu.Lock()
		defer mu.Unlock()

		// Unblock IP if block duration expired
		if unblockAt, blocked := blockedIPs[ip]; blocked {
			if time.Now().After(unblockAt) {
				delete(blockedIPs, ip)
			} else {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests. Try again later.",
				})
			}
		}

		// Update access count
		entry, exists := rateLimitMap[ip]
		if !exists || time.Since(entry.lastVisit) > rateWindow {
			rateLimitMap[ip] = &ipEntry{count: 1, lastVisit: time.Now()}
		} else {
			entry.count++
			entry.lastVisit = time.Now()

			if entry.count > maxRequestsPerWindow {
				blockedIPs[ip] = time.Now().Add(blockDuration)
				delete(rateLimitMap, ip)
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests. You are temporarily blocked.",
				})
			}
		}

		return next(c)
	}
}
