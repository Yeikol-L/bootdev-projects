package api

import (
	"net/http"

	"github.com/yeikol-l/bootdev-projects/chirpy/internal/helpers"
)

func (cfg *Handler) resetHandler(w http.ResponseWriter, r *http.Request) {
	if !cfg.isDev {
		helpers.RespondWithError(w, http.StatusForbidden, "")
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "Delete executed"})
}
