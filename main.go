package main

import (
	"database/sql"
	"log"
	"myserver/internal/database" 
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.Header().Set("X-My-Custom-Header", "This is a test header")

	w.WriteHeader(http.StatusOK) 

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbq            database.Queries
	platform       string
	secretk        string
	PolkaAPIKey    string
}

func main() {
	godotenv.Load()
	polkaKey := os.Getenv("POLKA_KEY")
	secretkey := os.Getenv("SECRETKEY")
	platform := os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Error connecting to db: ", err)
	}
	dbQueries := database.New(db)

	var Apicfg apiConfig
	Apicfg.dbq = *dbQueries
	Apicfg.platform = platform
	Apicfg.secretk = secretkey
	Apicfg.PolkaAPIKey = polkaKey

	ServeMux := http.NewServeMux()

	ServeMux.Handle("/app/", http.StripPrefix("/app", Apicfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	ServeMux.HandleFunc("GET /api/healthz", HomeHandler)

	ServeMux.HandleFunc("GET /admin/metrics", Apicfg.HitsHandler)
	ServeMux.HandleFunc("POST /admin/reset", Apicfg.ResetHandler)

	ServeMux.HandleFunc("POST /api/users", Apicfg.CreateUserHandler)
	ServeMux.HandleFunc("PUT /api/users", Apicfg.UpdateUserHandler)
	ServeMux.HandleFunc("POST /api/login", Apicfg.LoginHandler)
	ServeMux.HandleFunc("POST /api/refresh", Apicfg.RefreshHandler)
	ServeMux.HandleFunc("POST /api/revoke", Apicfg.RevokeHandler)

	ServeMux.HandleFunc("GET /api/chirps", Apicfg.GetAllChirpsHandler)
	ServeMux.HandleFunc("GET /api/chirps/{chirpID}", Apicfg.GetChirpHandler)
	ServeMux.HandleFunc("POST /api/chirps", Apicfg.ValidateChirpHandler)
	ServeMux.HandleFunc("DELETE /api/chirps/{chirpID}", Apicfg.DeleteChirpHandler)

	ServeMux.HandleFunc("POST /api/polka/webhooks", Apicfg.PolkaHandler)


	//swagger documentation:http://localhost:8080/swagger/ 
	// Настройка маршрута для Swagger UI
	ServeMux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger-ui/"))))

	// Настройка маршрута для swagger.json
	ServeMux.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	var MyServer http.Server
	MyServer.Handler = ServeMux
	MyServer.Addr = ":8080"

	log.Fatal(MyServer.ListenAndServe())
}

// на линуксе: go build -o out && ./out
// на windows: go build -o out.exe && out.exe
// http://localhost:8080

//dsn: psql postgres://postgres:password@localhost:5432/chirpy?sslmode=disable

//goose postgres postgres://postgres:password@localhost:5432/chirpy?sslmode=disable up migration_num
