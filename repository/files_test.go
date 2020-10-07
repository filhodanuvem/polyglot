package repository

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGetFiles(t *testing.T) {
	l := log.New()

	files := GetFiles("../testdata", l)

	if files == nil {
		t.Errorf("Should not return nil if the passed a valid files path")
	}
}
