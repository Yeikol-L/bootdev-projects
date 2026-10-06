package api

import (
	"net/http"
	"sync/atomic"

	"github.com/yeikol-l/bootdev-projects/chirpy/internal/database"
)

type Handler struct {
	db        *database.Queries
	isDev     bool
	counter   *atomic.Int32
	jwtSecret string
	polkaKey  string
}

func New(db *database.Queries, dev bool, jwtSecret string, polkaKey string) *Handler {
	return &Handler{
		db:        db,
		isDev:     dev,
		counter:   &atomic.Int32{},
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/login", h.loginHandler)
	mux.HandleFunc("POST /api/refresh", h.refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", h.revokeTokenHandler)
	mux.HandleFunc("POST /api/users", h.createUserHandler)
	mux.HandleFunc("PUT /api/users", h.updateUserHandler)
	mux.HandleFunc("POST /admin/reset", h.resetHandler)
	mux.HandleFunc("GET /admin/metrics", h.hitsMetricHandler)
	mux.HandleFunc("POST /api/chirps", h.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", h.getAllChirpsHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", h.deleteChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", h.getChirpHandler)

	mux.HandleFunc("POST /api/polka/webhooks", h.polkaWebhooksHandler)

	mux.HandleFunc("POST /api/validate_chirp", h.chirpHandler)
	mux.HandleFunc("GET /api/healthz", h.healthHandler)
	mux.Handle("/app/", h.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	return mux
}
