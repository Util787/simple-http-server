package main

import (
	"context"
	"log"
	"myserver/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Println("error getting token string: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuidFromToken, err := auth.ValidateJWT(tokenString, cfg.secretk)
	if err != nil {
		log.Println("error validating jwt: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpUuid, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Println("error parsing chirp id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp,err:=cfg.dbq.GetChirpsByChirpId(context.Background(),chirpUuid)
	if err!=nil{
		log.Println("error getting chirp by id: ",err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userUuidFromToken:=uuid.NullUUID{UUID:uuidFromToken,Valid:true}
	if chirp.UserID!=userUuidFromToken{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.dbq.DeleteChirpById(context.Background(), chirpUuid)
	if err != nil {
		log.Println("error deleting chirp by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
