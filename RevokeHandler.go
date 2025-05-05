package main

import (
	"context"
	"database/sql"
	"log"
	"myserver/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) RevokeHandler(w http.ResponseWriter, req *http.Request) {
	val := req.Header["Authorization"][0]
	refreshToken := val[7:]
	revokeParams := database.RevokeRefreshTokenParams{
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		Token:     refreshToken,
	}
	err := cfg.dbq.RevokeRefreshToken(context.Background(), revokeParams)
	if err != nil {
		log.Println("error revoking refresh token: ", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
