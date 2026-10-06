package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/auth"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/helpers"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (c *Handler) polkaWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if apiKey != c.polkaKey {
		helpers.RespondWithError(w, http.StatusUnauthorized, "ApiKey invalida")
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var polkaWebhookRequest PolkaWebhookRequest
	err = json.Unmarshal(bodyBytes, &polkaWebhookRequest)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if polkaWebhookRequest.Event != "user.upgraded" {
		helpers.RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	_, err = c.db.UpgradeUser(r.Context(), polkaWebhookRequest.Data.UserID)
	if err != nil {
		helpers.RespondWithJSON(w, http.StatusNotFound, nil)
		return
	}
	helpers.RespondWithJSON(w, http.StatusNoContent, nil)

}
