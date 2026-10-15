package httpx

import (
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request, maxHops int, trusted []*net.IPNet) string {
	remote, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remote = r.RemoteAddr
	}

	// Si la conexión no viene de un proxy confiable, ignoramos el header
	if maxHops <= 0 || !isTrusted(remote, trusted) {
		return remote
	}

	var parts []string
	for _, v := range r.Header.Values("X-Forwarded-For") {
		parts = append(parts, strings.Split(v, ",")...)
	}

	ip := remote
	for i, hops := len(parts)-1, 0; i >= 0 && hops < maxHops; i, hops = i-1, hops+1 {
		candidate := strings.TrimSpace(parts[i])
		if net.ParseIP(candidate) == nil {
			break
		}
		ip = candidate
		if !isTrusted(candidate, trusted) {
			break // primera IP que no es proxy = cliente real
		}
	}
	return ip
}

func isTrusted(ip string, trusted []*net.IPNet) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, n := range trusted {
		if n.Contains(parsed) {
			return true
		}
	}
	return false
}
