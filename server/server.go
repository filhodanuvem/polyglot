package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
	"github.com/filhodanuvem/polyglot/stats"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Languages []repository.Counter `json:"languages"`
	Username  string               `json:"user"`
}

type Config struct {
	Port     string
	Host     string
	TempPath string
	Log      *logrus.Logger
}

func getLanguages(w http.ResponseWriter, req *http.Request, config Config) {

	l := config.Log

	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method Not Allowed"}`))
		l.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	query := req.URL.Query()
	username := query.Get("user")
	queryLimit := query.Get("limit")

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

	repos, err := github.GetRepositories(username)
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

	stats := stats.GetStatisticsAsync(config.TempPath, repos, config.Log)
	firstLanguages := stats.FirstLanguages(int(limit))

	response := &Response{
		Languages: firstLanguages,
		Username:  username,
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

func Serve(config Config) {

	l := config.Log

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		getLanguages(w, req, config)
	})
	serverAddress := config.Host + ":" + config.Port

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
	serverErr := http.ListenAndServe(serverAddress, nil)
	if serverErr != nil {
		l.Error(serverErr)
	}
}
