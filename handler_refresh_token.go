package main

import (
	"net/http"
	"time"

	"github.com/cadimodev/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	refreshTokenResult, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if time.Now().UTC().After(refreshTokenResult.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if refreshTokenResult.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenResult.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}

	expirationTime := time.Hour

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
