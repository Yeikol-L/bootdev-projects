//lint:file-ignore U1000 Some code in this file is used later

package httpserver

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bootdotdev/learn-web-security/internal/httpx"
	"github.com/bootdotdev/learn-web-security/internal/logging"
	"github.com/bootdotdev/learn-web-security/internal/templates"
)

type middleware func(http.Handler) http.Handler

func applyMiddleware(handler http.Handler, middlewareChain ...middleware) http.Handler {
	for _, currentMiddleware := range slices.Backward(middlewareChain) {
		handler = currentMiddleware(handler)
	}
	return handler
}

var OriginHeader = "Origin"
var OriginHeaderFallback = "Referer"

func sourceValidation(appOrigin string, renderer *templates.Renderer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				next.ServeHTTP(w, r)
				return
			}
			originHeaderValue := r.Header.Get(OriginHeader)
			refererHeaderValue := r.Header.Get(OriginHeaderFallback)
			requestOrigin := originHeaderValue
			if originHeaderValue == "" {
				if url, err := url.Parse(refererHeaderValue); err == nil {
					requestOrigin = fmt.Sprintf("%s://%s", url.Scheme, url.Host)
				}
			}
			if requestOrigin != appOrigin {
				httpx.RespondWithErrorPage(w, renderer, http.StatusForbidden, "Not allowed", "Please use the oficial domain")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func productCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" && request.URL.Path == "/api/products" {
			responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		}
		next.ServeHTTP(responseWriter, request)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonceBytes := make([]byte, 16)
		if _, err := rand.Read(nonceBytes); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		nonce := base64.StdEncoding.EncodeToString(nonceBytes)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'self' 'nonce-%s'; script-src 'self' 'nonce-%s'; style-src 'self'; img-src 'self' data:; frame-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'self'", nonce, nonce))
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Origin-Agent-Cluster", "?1")
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("X-XSS-Protection", "0")

		r = r.WithContext(httpx.WithCSPNonce(r.Context(), nonce))
		next.ServeHTTP(w, r)
	})
}

func overrideCrosOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		next.ServeHTTP(w, r)
	})
}

func recoverPanics(logger *logging.Logger, renderer *templates.Renderer) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					_ = logger.Event("unhandled_error", map[string]any{
						"method":  request.Method,
						"path":    request.URL.Path,
						"message": fmt.Sprint(recovered),
					})
					if err := httpx.RespondWithErrorPage(responseWriter, renderer, http.StatusInternalServerError, "Unhandled Error", fmt.Sprint(recovered)); err != nil {
						http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(responseWriter, request)
		})
	}
}

func LoadShedder(cap int, delay int) func(http.Handler) http.Handler {
	if cap <= 0 || delay <= 0 {
		panic("Cap or delay cant be negative")
	}

	limit := strconv.Itoa(cap)
	retry := strconv.Itoa(delay)
	return func(next http.Handler) http.Handler {
		bufferChan := make(chan struct{}, cap)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-In-Flight-Limit", limit)
			select {
			case bufferChan <- struct{}{}:
				defer func() { <-bufferChan }()
				next.ServeHTTP(w, r)
			default:
				w.Header().Add("Retry-After", retry)
				httpx.RespondWithJSON(w, 503, map[string]any{"error": "Service is at capacity"})
			}
		})
	}
}

func SearchThrottle(renderer *templates.Renderer) func(http.Handler) http.Handler {
	opts := rateLimitOptions{
		window:  time.Second,
		maximum: 5,
		key: func(r *http.Request) string {
			return "SEARCH-CONSTAT"
		},
		onLimit: func(w http.ResponseWriter, r *http.Request, rls rateLimitState) {
			// w.Header().Set("Retry-After", strconv.FormatInt(int64(rls.retryAfterSeconds), 10))
			if renderErr := httpx.RespondWithErrorPage(w, renderer, http.StatusTooManyRequests, "Search Is Busy", "Try again shortly."); renderErr != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		},
		now: func() time.Time { return time.Now() },
	}
	return fixedWindowRateLimiter(opts)
}

type rateLimitCounter struct {
	count   int
	resetAt time.Time
}

type rateLimitOptions struct {
	window  time.Duration
	maximum int
	key     func(*http.Request) string
	onLimit func(http.ResponseWriter, *http.Request, rateLimitState)
	now     func() time.Time
}

type rateLimitState struct {
	limit             int
	remaining         int
	resetAt           time.Time
	retryAfterSeconds int
}

type fixedWindowLimiter struct {
	options       rateLimitOptions
	counters      map[string]rateLimitCounter
	countersMutex sync.Mutex
	nextSweepAt   time.Time
}

func newFixedWindowLimiter(options rateLimitOptions) *fixedWindowLimiter {
	validateRateLimitOptions(options)
	if options.now == nil {
		options.now = time.Now
	}
	if options.key == nil {
		options.key = clientIPKey
	}
	return &fixedWindowLimiter{
		options:     options,
		counters:    make(map[string]rateLimitCounter),
		nextSweepAt: options.now().Add(options.window),
	}
}

func validateRateLimitOptions(options rateLimitOptions) {
	if options.window <= 0 {
		panic("rate-limit window must be positive")
	}
	if options.maximum <= 0 {
		panic("rate-limit maximum must be positive")
	}
}

func (limiter *fixedWindowLimiter) consume(request *http.Request) (rateLimitState, bool) {
	now := limiter.options.now()
	key := limiter.options.key(request)

	limiter.countersMutex.Lock()
	defer limiter.countersMutex.Unlock()
	limiter.sweepExpiredCounters(now)

	counter, exists := limiter.counters[key]
	if !exists || !now.Before(counter.resetAt) {
		counter = rateLimitCounter{resetAt: now.Add(limiter.options.window)}
	}
	retryAfterSeconds := max(1, int(math.Ceil(counter.resetAt.Sub(now).Seconds())))
	if counter.count >= limiter.options.maximum {
		return rateLimitState{
			limit:             limiter.options.maximum,
			remaining:         0,
			resetAt:           counter.resetAt,
			retryAfterSeconds: retryAfterSeconds,
		}, true
	}

	counter.count++
	limiter.counters[key] = counter
	return rateLimitState{
		limit:             limiter.options.maximum,
		remaining:         limiter.options.maximum - counter.count,
		resetAt:           counter.resetAt,
		retryAfterSeconds: retryAfterSeconds,
	}, false
}

func (limiter *fixedWindowLimiter) sweepExpiredCounters(now time.Time) {
	if now.Before(limiter.nextSweepAt) {
		return
	}
	for counterKey, counter := range limiter.counters {
		if !now.Before(counter.resetAt) {
			delete(limiter.counters, counterKey)
		}
	}
	limiter.nextSweepAt = now.Add(limiter.options.window)
}

func (limiter *fixedWindowLimiter) reject(responseWriter http.ResponseWriter, request *http.Request, state rateLimitState) {
	setRateLimitHeaders(responseWriter, state)
	responseWriter.Header().Set("Retry-After", strconv.Itoa(state.retryAfterSeconds))
	if limiter.options.onLimit != nil {
		limiter.options.onLimit(responseWriter, request, state)
		return
	}
	httpx.RespondWithJSON(responseWriter, http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
}

func fixedWindowRateLimiter(options rateLimitOptions) middleware {
	// validateRateLimitOptions(options)
	limiter := newFixedWindowLimiter(options)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				state, limited := limiter.consume(r)
				if limited {
					limiter.reject(w, r, state)
					return
				}
				setRateLimitHeaders(w, state)
				next.ServeHTTP(w, r)
			})
	}
}

func setRateLimitHeaders(responseWriter http.ResponseWriter, state rateLimitState) {
	responseWriter.Header().Set("RateLimit-Limit", strconv.Itoa(state.limit))
	responseWriter.Header().Set("RateLimit-Remaining", strconv.Itoa(state.remaining))
	responseWriter.Header().Set("RateLimit-Reset", strconv.FormatInt(state.resetAt.Unix()+boolToInt64(state.resetAt.Nanosecond() > 0), 10))
}

func clientIPKey(request *http.Request) string {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	if request.RemoteAddr != "" {
		return request.RemoteAddr
	}
	return "unknown"
}

func clientIPKeyWithTrustedProxies(trustedProxyHops int) func(*http.Request) string {
	return func(request *http.Request) string {
		if trustedProxyHops > 0 {
			forwardedAddresses := strings.Split(request.Header.Get("X-Forwarded-For"), ",")
			selectedIndex := len(forwardedAddresses) - trustedProxyHops
			if selectedIndex >= 0 && selectedIndex < len(forwardedAddresses) {
				if selectedAddress := strings.TrimSpace(forwardedAddresses[selectedIndex]); selectedAddress != "" {
					return selectedAddress
				}
			}
		}
		return clientIPKey(request)
	}
}

func boolToInt64(value bool) int64 {
	if value {
		return 1
	}
	return 0
}
