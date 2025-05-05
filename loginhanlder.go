package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"myserver/internal/auth"
	"myserver/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginReqBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRespBody struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var LoginReq LoginReqBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&LoginReq)
	if err != nil {
		log.Println("error decoding login request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserByEmail, err := cfg.dbq.GetUserByEmail(context.Background(), LoginReq.Email)
	if err != nil {
		log.Println("error getting user from db: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	password := LoginReq.Password

	err = auth.CheckPasswordHash([]byte(password), []byte(UserByEmail.HashedPassword))
	if err != nil {
		log.Println("error checking password: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accesTokenExpireTime := time.Hour
	log.Println("expiretime: ", accesTokenExpireTime) // to delete

	accesToken, err := auth.MakeJWT(UserByEmail.ID, cfg.secretk, accesTokenExpireTime)
	if err != nil {
		log.Println("error making jwt: ", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("error making refresh token: ", err)
		return
	}

	refreshTokenExpireDate := time.Now().Add(time.Hour * 1440) // expires after 60 days since created

	RefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: UserByEmail.ID, Valid: true},
		ExpiresAt: refreshTokenExpireDate,
		RevokedAt: sql.NullTime{Valid: false},
	}
	_, err = cfg.dbq.CreateRefreshToken(context.Background(), RefreshTokenParams)
	if err != nil {
		log.Println("error creating refresh token: ", err)
		return
	}

	response := LoginRespBody{
		ID:           UserByEmail.ID,
		CreatedAt:    UserByEmail.CreatedAt,
		UpdatedAt:    UserByEmail.UpdatedAt,
		Email:        UserByEmail.Email,
		IsChirpyRed: UserByEmail.IsChirpyRed,
		Token:        accesToken,
		RefreshToken: refreshToken,
	}

	respdata, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling on server side: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respdata)
}
