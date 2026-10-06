package api

import "net/http"

func (cfg *Handler) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.counter.Add(1)
		next.ServeHTTP(w, r)
	})
}
