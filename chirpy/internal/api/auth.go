package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/auth"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/database"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/helpers"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func (c *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var body LoginRequest
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = validate.Struct(body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := c.db.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "El usuario o contraseña no son validos")
		return
	}
	passwordMatch, err := auth.CheckPassword(body.Password, user.HashedPassword)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !passwordMatch {
		helpers.RespondWithError(w, http.StatusUnauthorized, "El usuario o contraseña no son validos")
		return
	}

	jwt, err := auth.MakeJWT(user.ID, c.jwtSecret, time.Hour)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}
	refreshToken, err := c.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, LoginResponse{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        jwt,
		IsChirpyRed:  user.IsChirpyRed,
		RefreshToken: refreshToken.Token,
	})

}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func (c *Handler) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	tokenData, err := c.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token no existe")
		return
	}
	if tokenData.RevokedAt.Valid {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token revocado")
		return
	}
	if tokenData.ExpiresAt.Before(time.Now()) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token expirado")
		return
	}
	jwt, err := auth.MakeJWT(tokenData.UserID, c.jwtSecret, time.Hour)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, RefreshTokenResponse{Token: jwt})
}
func (c *Handler) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	_, err = c.db.RevokeToken(r.Context(), token)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Token no existe")
		return
	}
	helpers.RespondWithJSON(w, http.StatusNoContent, nil)
}
