package main

import (
	"context"
	"encoding/json"
	"log"
	"myserver/internal/auth"
	"myserver/internal/database"
	"net/http"
	"time"
)

type UpdateUserReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
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

	var DecodedReqBody UpdateUserReqBody

	decoder := json.NewDecoder(req.Body)

	decoder.Decode(&DecodedReqBody)

	hashedPassword, err := auth.HashPassword([]byte(DecodedReqBody.Password))
	if err != nil {
		log.Println("error hashing password: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var params database.UpdateUserParams
	params.Email = DecodedReqBody.Email
	params.ID = uuidFromToken
	params.HashedPassword = string(hashedPassword)
	params.UpdatedAt = time.Now()

	err = cfg.dbq.UpdateUser(context.Background(), params)
	if err != nil {
		log.Println("error updating user: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
