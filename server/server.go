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

func getLanguages(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method Not Allowed"}`))
		fmt.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusMethodNotAllowed)
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
		fmt.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusUnprocessableEntity)
		return
	}

	repos, err := github.GetRepositories(username)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	l := logrus.New()

	stats := stats.GetStatisticsAsync("/tmp/polyglot", repos, l)
	firstLanguages := stats.FirstLanguages(int(limit))

	response := &Response{
		Languages: firstLanguages,
		Username:  username,
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	fmt.Printf("%v - %v - %v \n", req.Method, req.URL, http.StatusOK)
}

func Serve(host string, port string) {

	http.HandleFunc("/", getLanguages)
	serverAddress := host + ":" + port

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(serverErr.Error())
	}
}
