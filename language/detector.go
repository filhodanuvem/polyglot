package language

import (
	"io/ioutil"

	"github.com/go-enry/go-enry/v2"
)

func DetectByFile(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	lang := enry.GetLanguage(file, content)

	return lang, nil
}
