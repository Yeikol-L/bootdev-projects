package api

import (
	"encoding/json"
	"fmt"

	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/auth"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/database"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/helpers"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsChirpyRead bool      `json:"is_chirpy_red"`
}

func (cfg *Handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Could not parse the body")
		return
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not process the request")
		return
	}
	createUserParams := database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
	}
	newUser, err := cfg.db.CreateUser(r.Context(), createUserParams)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, UserResponse{
		ID:           newUser.ID,
		Email:        newUser.Email,
		CreatedAt:    newUser.CreatedAt,
		UpdatedAt:    newUser.UpdatedAt,
		IsChirpyRead: newUser.IsChirpyRed,
	})
}

func (c *Handler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var updateUserRequest UpdateUserRequest
	err = json.Unmarshal(bodyBytes, &updateUserRequest)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = c.db.GetUserById(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	hashedPassword, err := auth.HashPassword(updateUserRequest.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updatedUser, err := c.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          updateUserRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, UserResponse{
		ID:           updatedUser.ID,
		Email:        updateUserRequest.Email,
		CreatedAt:    updatedUser.CreatedAt,
		UpdatedAt:    updatedUser.UpdatedAt,
		IsChirpyRead: updatedUser.IsChirpyRed,
	})
}
