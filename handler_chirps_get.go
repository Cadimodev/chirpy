package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetChirpsOrderByCreatedAt(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	sortString := r.URL.Query().Get("sort")

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {

		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	switch sortString {
	case "asc":
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	case "desc":
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
