package repository

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	l, hook := test.NewNullLogger()

	files := GetFiles("../testdata", l)

	if files == nil {
		t.Errorf("Should not return nil if the passed a valid files path")
	}

	// Test the GetFiles logs output
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "Listing files of repo", hook.LastEntry().Message)

	files = GetFiles("", l)

	if files != nil {
		t.Errorf("Should return nil if the path leads to a invalid files folder")
	}
}
