package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// GetChirpHandler finding chirp by chirpID
//	@Summary		finding chirp by chirpID
//	@Description	finding chirp by chirpID provided in path
//	@Tags			Chirps
//	@Produce		json
//	@Param			chirpID	path		string	true	"chirpID"
//	@Success		200		{object}	ChirpValidRespBody
//	@Router			/api/chirps/{chirpID} [get]
func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, req *http.Request) {

	chirpid, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Println("error parsing chirp id: ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dbchirp, err := cfg.dbq.GetChirpsByChirpId(context.Background(), chirpid)
	if err != nil {
		log.Println("error getting chirp from db: ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chirp := ChirpValidRespBody{
		ID:        dbchirp.ID.String(),
		CreatedAt: dbchirp.CreatedAt,
		UpdatedAt: dbchirp.UpdatedAt,
		Body:      dbchirp.Body,
		UserID:    dbchirp.UserID.UUID.String(),
	}

	respdata, err := json.Marshal(chirp)
	if err != nil {
		log.Println("error marshalling chirp: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respdata)
}
