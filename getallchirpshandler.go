package main

import (
	"context"
	"encoding/json"
	"log"
	"myserver/internal/database"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

// GetAllChirpsHandler getting all Chirps or all Chirps by the author
//	@Summary		Getting all Chirps or all Chirps by the author
//	@Description	Getting all Chirps if query string is empty or all Chirps by the author if authorID is provided in query string, you also can provide sort order in query string
//	@Tags			Chirps
//	@Produce		json
//	@Param			author_id	query	string	false	"authorID"
//	@Param			sort		query	string	false	"Sort order, can either be 'asc' or 'desc', asc by default"
//	@Success		200			{array}	ChirpValidRespBody
//	@Router			/api/chirps [get]
func (cfg *apiConfig) GetAllChirpsHandler(w http.ResponseWriter, req *http.Request) {

	// authorIDQueryParam is a string that contains the value of the author_id query parameter if it exists, or an empty string if it doesn't
	// For example authorIDQueryParam = "1" when:
	// GET http://localhost:8080/api/chirps?author_id=1
	authorIDQueryParam := req.URL.Query().Get("author_id")

	chirps, err := cfg.GetChirps(authorIDQueryParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//sort query expecting "desc" or "asc"
	SortType := req.URL.Query().Get("sort")
	if SortType == "desc" {
		slices.Reverse(chirps)
	}

	ChirpsResp := make([]ChirpValidRespBody, len(chirps))

	for i, chirp := range chirps {
		ChirpsResp[i] = ChirpValidRespBody{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.UUID.String(),
		}
	}

	respdata, err := json.Marshal(ChirpsResp)
	if err != nil {
		log.Println("error marshalling on server side")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respdata)
}

func (cfg *apiConfig) GetChirps(IdFromQueryString string) ([]database.Chirp, error) {
	if IdFromQueryString != "" {

		UserID, err := uuid.Parse(IdFromQueryString)
		if err != nil {
			log.Println("error parsing uuid from query string: ", err)
			return nil, err
		}

		chirps, err := cfg.dbq.GetChirpsByUserId(context.Background(), uuid.NullUUID{UUID: UserID, Valid: true})
		if err != nil {
			log.Println("error getting chirps from db: ", err)
			return nil, err
		}

		return chirps, nil
	}

	chirps, err := cfg.dbq.GetAllChirps(context.Background())
	if err != nil {
		log.Println("error getting chirps from db: ", err)
		return nil, err
	}

	return chirps, nil
}
