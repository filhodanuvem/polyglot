package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	Items []struct {
		URL string `json:"html_url"`
	}
}

func GetRepositories(username string) ([]string, error) {
	var repos []string
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/search/repositories?q=user:%s", username))
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
