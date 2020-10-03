package language

import "testing"

func TestDetectByFile(t *testing.T) {
	_, err := DetectByFile("")

	if err == nil {
		t.Errorf("Should trow an error if the file is invalid")
	}
}
