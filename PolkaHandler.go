package main

import (
	"context"
	"encoding/json"
	"log"
	"myserver/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

type PolkaReqBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) PolkaHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	ReqApiKey,err:=auth.GetAPIKey(req.Header)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cfg.PolkaAPIKey!=ReqApiKey{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder:=json.NewDecoder(req.Body)

	var ReqBody PolkaReqBody

	err=decoder.Decode(&ReqBody)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error decoding Polka request body: ",err)
		return
	}

	if ReqBody.Event!="user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return
	}

	exists,err:=cfg.dbq.UserExists(context.Background(),ReqBody.Data.UserID)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = cfg.dbq.AddRedSubscription(context.Background(),ReqBody.Data.UserID)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}