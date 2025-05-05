package main

import (
	"context"
	"encoding/json"
	"log"
	"myserver/internal/auth"
	"myserver/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var BannedWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type ChirpReqBody struct { // this is what func gets in the POST request
	Body string `json:"body"`
}

type ChirpErrorRespBody struct {
	Error string `json:"error"`
}

type ChirpValidRespBody struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    string `json:"user_id"`
}

func (cfg *apiConfig) ValidateChirpHandler(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close() // didnt find info about if this line is necessary but just in case

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

	decoder := json.NewDecoder(req.Body)

	var reqbody ChirpReqBody

	err = decoder.Decode(&reqbody)
	if err != nil {
		var ErrResp ChirpErrorRespBody
		ErrResp.Error = err.Error()

		respdata, ServErr := json.Marshal(ErrResp)
		if ServErr != nil {
			log.Println("error marshalling on server side: ", ServErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respdata)
		return
	}

	if len(reqbody.Body) > 140 {
		var ErrResp ChirpErrorRespBody
		ErrResp.Error = "chirp is too long"

		respdata, ServErr := json.Marshal(ErrResp)
		if ServErr != nil {
			log.Println("error marshalling on server side: ", ServErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respdata)
		return
	}

	Message := strings.Fields(reqbody.Body)
	for i, word := range Message {
		if _, exist := BannedWords[strings.ToLower(word)]; exist {
			Message[i] = "****"
		}
	}

	cleanedbody := strings.Join(Message, " ")

	var NewChirpParams database.CreateChirpParams
	NewChirpParams.ID = uuid.New()
	NewChirpParams.CreatedAt = time.Now()
	NewChirpParams.UpdatedAt = time.Now()
	NewChirpParams.Body = cleanedbody
	NewChirpParams.UserID = uuid.NullUUID{UUID: uuidFromToken, Valid: true}
	_, err = cfg.dbq.CreateChirp(context.Background(), NewChirpParams)
	if err != nil {
		log.Println("error creating chirp in db: ", err)
		var ErrResp ChirpErrorRespBody
		ErrResp.Error = "error creating chirp in db"
		w.WriteHeader(http.StatusBadRequest)
		respdata, err := json.Marshal(ErrResp)
		if err != nil {
			log.Println("error marshaling on server side: ", err)
			return
		}
		w.Write(respdata)
		return
	}

	var ValidResp ChirpValidRespBody
	ValidResp.ID = NewChirpParams.ID.String()
	ValidResp.CreatedAt = NewChirpParams.CreatedAt
	ValidResp.UpdatedAt = NewChirpParams.UpdatedAt
	ValidResp.Body = NewChirpParams.Body
	ValidResp.UserID = NewChirpParams.UserID.UUID.String()

	respdata, ServErr := json.Marshal(ValidResp)
	if ServErr != nil {
		log.Println("error marshalling on server side")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respdata)
}
