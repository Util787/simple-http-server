package main

import (
	"context"
	"encoding/json"
	"log"
	"myserver/internal/auth"
	"net/http"
	"time"
)

type RefreshValidRespBody struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) RefreshHandler(w http.ResponseWriter, req *http.Request) {
	val := req.Header["Authorization"][0]
	refreshToken := val[7:]
	dbRefreshToken, err := cfg.dbq.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		log.Println("error getting refresh token from db: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	timeNowWithoutUTC:=time.Now().Add(time.Hour*3)

	if dbRefreshToken.ExpiresAt.Before(timeNowWithoutUTC.UTC()) {
		log.Println("expired token")
		log.Printf("Token ExpiresAt: %v, Current Time: %v", dbRefreshToken.ExpiresAt, timeNowWithoutUTC.UTC())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("token is not expired")
	log.Printf("Token ExpiresAt: %v, Current Time: %v", dbRefreshToken.ExpiresAt, timeNowWithoutUTC.UTC())
	log.Println("IN UTC: ",dbRefreshToken.ExpiresAt.UTC(),time.Now().UTC())

	if !dbRefreshToken.RevokedAt.Time.IsZero() {
		log.Println("revoked token")
		log.Printf("Token RevokedAt: %v, Current Time: %v", dbRefreshToken.RevokedAt.Time, timeNowWithoutUTC.UTC())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("token is not revoked")
	log.Printf("Token RevokedAt: %v, Current Time: %v", dbRefreshToken.RevokedAt.Time, timeNowWithoutUTC.UTC())

	accesTokenExpireTime := time.Hour
	accessToken, err := auth.MakeJWT(dbRefreshToken.UserID.UUID, cfg.secretk, accesTokenExpireTime)
	if err != nil {
		log.Println("error making jwt in refresh handler: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := RefreshValidRespBody{
		Token: accessToken,
	}

	respdata, err := json.Marshal(resp)
	if err != nil {
		log.Println("error marshalling response data: ", respdata)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respdata)
}
