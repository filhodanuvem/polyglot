package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/filhodanuvem/polyglot/source"
)

type response struct {
	URL           string `json:"web_url"`
	DefaultBranch string `json:"default_branch"`
}

func (r response) GetUrl() string {
	return r.URL
}

func GetRepositories(username string) ([]source.ProviderRepo, error) {
	repos := []source.ProviderRepo{}
	resp, err := http.Get(fmt.Sprintf("https://gitlab.com/api/v4/users/%s/projects", username))
	if err != nil {
		return repos, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make([]response, 0)
	json.Unmarshal(body, &jsonResponse)

	for _, repo := range jsonResponse {
		repos = append(repos, source.ProviderRepo{
			URL:           repo.URL,
			DefaultBranch: repo.DefaultBranch,
		})
	}

	return repos, nil
}
