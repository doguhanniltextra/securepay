package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// client struct to track request count and window start time for each IP.
type client struct {
	count     int
	startTime time.Time
}

var (
	// clients map stores the rate limit state for each IP.
	// In a production environment with multiple replicas, this should be in Redis.
	// For this task, strictly in-memory map is used as requested.
	clients = make(map[string]*client)
	mu      sync.Mutex
)

// RateLimitMiddleware enforces a limit of 100 requests per minute per IP.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		mu.Lock()
		c, exists := clients[ip]
		if !exists {
			clients[ip] = &client{count: 1, startTime: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Check if the 1-minute window has passed
		if time.Since(c.startTime) > time.Minute {
			// Reset counter and start new window
			c.count = 1
			c.startTime = time.Now()
		} else {
			// Increment counter
			c.count++
		}

		// Check limit
		if c.count > 100 {
			mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client's real IP address from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (standard for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}

	// Check X-Real-Ip header
	realIP := r.Header.Get("X-Real-Ip")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
