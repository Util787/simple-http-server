package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"myserver/internal/auth"
	"myserver/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type CreateUserReqBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) CreateUserHandler(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)

	reqbody := CreateUserReqBody{}

	err := decoder.Decode(&reqbody)
	if err != nil {
		log.Println("error decoding: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := reqbody.Email
	password := reqbody.Password
	hashed_password, err := auth.HashPassword([]byte(password))
	if err != nil {
		log.Println("error hashing password: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := database.CreateUserParams{}
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	params.Email = email
	params.HashedPassword = string(hashed_password)

	dbUser, err := cfg.dbq.CreateUser(context.Background(), params)
	if err != nil {
		log.Println(errors.New("error creating user in db: "), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// mapping it from database.User to User struct to let me control the JSON keys
	CreatedUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	okRespData, err := json.Marshal(CreatedUser)
	if err != nil {
		log.Println("err marshalling createduser: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(okRespData)
}
