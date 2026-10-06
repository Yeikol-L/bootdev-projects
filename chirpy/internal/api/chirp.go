package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/auth"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/database"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/helpers"
)

var BANNED_WORDS = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

type CreateChirpRequest struct {
	Body string `json:"body"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func sanitizeChirpBody(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}
	words := strings.Split(body, " ")
	for i, w := range words {
		for _, bw := range BANNED_WORDS {
			if strings.EqualFold(w, bw) {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " "), nil
}

func (c *Handler) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Body string `json:"body"`
	}
	var body RequestBody
	data, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}
	chirpBody, err := sanitizeChirpBody(body.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
	}
	helpers.RespondWithJSON(w, 200, map[string]string{"cleaned_body": chirpBody})
}

func (c *Handler) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var reqBody CreateChirpRequest
	err = json.Unmarshal(reqBodyBytes, &reqBody)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirpBody, err := sanitizeChirpBody(reqBody.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := c.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: chirpBody, UserID: userId})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, ChirpResponse{ID: chirp.ID, Body: chirp.Body, UserID: chirp.UserID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt})
}

func (c *Handler) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.URL.Query().Get("author_id"))
	sort := r.URL.Query().Get("sort")
	if sort != "desc" && sort != "asc" {
		sort = "asc"
	}
	params := database.GetAllChirpsParams{
		UserID: uuid.NullUUID{UUID: userID, Valid: err == nil},
		Sort:   sort,
	}
	chirps, err := c.db.GetAllChirps(r.Context(), params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirpResponses := []ChirpResponse{}

	for _, v := range chirps {
		chirpResponses = append(chirpResponses, ChirpResponse{
			ID:        v.ID,
			Body:      v.Body,
			UserID:    v.UserID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	helpers.RespondWithJSON(w, http.StatusOK, chirpResponses)
}

func (c *Handler) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := c.db.GetChirpByID(r.Context(), id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}

	helpers.RespondWithJSON(w, http.StatusOK, chirpResponse)
}

func (c *Handler) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirpData, err := c.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Chirp no encontrado")
		return
	}
	if chirpData.UserID != userId {
		helpers.RespondWithError(w, http.StatusForbidden, "No tiene permisos suficientes")
		return
	}
	err = c.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusNoContent, nil)
}
