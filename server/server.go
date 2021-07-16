package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	_ "github.com/filhodanuvem/polyglot/docs"
	"github.com/filhodanuvem/polyglot/repository"
	"github.com/filhodanuvem/polyglot/source/github"
	"github.com/filhodanuvem/polyglot/stats"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Response struct {
	Languages []repository.Counter `json:"languages"`
	Username  string               `json:"user"`
	Debug     debugResponse
}

type debugResponse struct {
	TimeToGetRepositoriesMs int64
	TimeToGetStatisticsMs   int64
}

type Config struct {
	Port     string
	Host     string
	TempPath string
	Log      *logrus.Logger
}

// GetLanguages godoc
// @Summary Get languagens
// @Description Get languagens by user
// @ID get-languages-by=username
// @Accept  json
// @Produce  json
// @Param  username query string true "github username"
// @Success 200
// @Router /search [get]
func GetLanguages(w http.ResponseWriter, req *http.Request, config Config) {

	l := config.Log

	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method Not Allowed"}`))
		l.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	query := req.URL.Query()
	username := query.Get("username")
	queryLimit := query.Get("limit")

	provider := "github"

	provider = query.Get("provider")

	limit, err := strconv.ParseInt(queryLimit, 10, 64)

	if err != nil {
		limit = 5
	}

	if len(username) <= 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error": "Missing username!"}`))
		l.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusUnprocessableEntity)
		return
	}

	beforeGetRepositories := time.Now()
	repos, err := github.GetRepositories(username)
	diff := time.Now().Sub(beforeGetRepositories)
	timeToGetRepositoriesMs := diff.Milliseconds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "` + err.Error() + `"} `))
		l.Errorf("%v - %v - %v \n", req.Method, req.URL, http.StatusInternalServerError)
		return
	}

	if len(repos) < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "This user has no public repositories"} `))
		l.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusNotFound)
		return
	}

	beforeGetStatistics := time.Now()
	stats := stats.GetStatisticsAsync(config.TempPath, provider, repos, config.Log)
	diff = time.Now().Sub(beforeGetStatistics)
	timeToGetStatisticsMs := diff.Milliseconds()
	firstLanguages := stats.FirstLanguages(int(limit))

	response := &Response{
		Languages: firstLanguages,
		Username:  username,
		Debug: debugResponse{
			TimeToGetRepositoriesMs: timeToGetRepositoriesMs,
			TimeToGetStatisticsMs:   timeToGetStatisticsMs,
		},
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "` + err.Error() + `"} `))
		l.Errorf("%v - %v - %v \n", req.Method, req.URL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	l.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusOK)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func Serve(config Config) {

	l := config.Log
	serverAddress := config.Host + ":" + config.Port
	swaggerURL := fmt.Sprintf("http://%s/swagger/doc.json", serverAddress)
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerURL), //The url pointing to API definition"
	))
	fmt.Println(swaggerURL)

	r.Get("/search", func(w http.ResponseWriter, req *http.Request) {
		GetLanguages(w, req, config)
	})

	listener, err := net.Listen("tcp", ":"+config.Port)

	if err != nil {
		l.Error(err)
		return
	}
	listener.Close()

	fmt.Println("\033[31m             _             _       _")
	fmt.Println("\033[31m _ __   ___ | |_   _  __ _| | ___ | |_")
	fmt.Println("\033[33m| '_ \\ / _ \\| | | | |/ _` | |/ _ \\| __|")
	fmt.Println("\033[32m| |_) | (_) | | |_| | (_| | | (_) | |_")
	fmt.Println("\033[34m| .__/ \\___/|_|\\__, |\\__, |_|\\___/ \\__|")
	fmt.Println("\033[35m|_|            |___/ |___/\033[0m")
	fmt.Printf("\nServer started at http://%v\n\n", serverAddress)
	serverErr := http.ListenAndServe(":"+config.Port, r)
	if serverErr != nil {
		l.Error(serverErr)
	}
}
