package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close() // didnt find info about if this line is necessary but just in case

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	cfg.dbq.DeleteAllUsers(context.Background())
}
