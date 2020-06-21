package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type response struct {
	Items []struct {
		URL string `json:"html_url"`
	}
}

func GetRepositories(profile string) ([]string, error) {
	var repos []string
	resp, err := http.Get("https://api.github.com/search/repositories?q=user:filhodanuvem")
	if err != nil {
		return repos, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse response
	json.Unmarshal(body, &jsonResponse)

	for i := range jsonResponse.Items {
		repos = append(repos, jsonResponse.Items[i].URL)
	}

	return repos, nil
}
