package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	URL string `json:"web_url"`
}

func (r response) GetUrl() string {
	return r.URL
}

func GetRepositories(username string) ([]string, error) {
	var repos []string
	resp, err := http.Get(fmt.Sprintf("https://gitlab.com/api/v4/users/%s/projects", username))
	if err != nil {
		return repos, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make([]response, 0)
	json.Unmarshal(body, &jsonResponse)

	for _, repo := range jsonResponse {
		repos = append(repos, repo.URL)
	}

	return repos, nil
}
