package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/filhodanuvem/polyglot/source"
)

type response struct {
	Items []struct {
		URL           string `json:"html_url"`
		DefaultBranch string `json:"default_branch"`
	}
}

func GetRepositories(username string) ([]source.ProviderRepos, error) {
	repos := []source.ProviderRepos{}
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/search/repositories?q=user:%s", username))
	if err != nil {
		return repos, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse response
	json.Unmarshal(body, &jsonResponse)

	for i := range jsonResponse.Items {
		repos = append(repos, source.ProviderRepos{
			URL:           jsonResponse.Items[i].URL,
			DefaultBranch: jsonResponse.Items[i].DefaultBranch,
		})
	}

	return repos, nil
}
