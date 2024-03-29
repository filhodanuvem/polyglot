package language

import "testing"

func TestDetectByFile(t *testing.T) {
	_, err := DetectByFile("")

	if err == nil {
		t.Errorf("Should trow an error if the file is invalid")
	}

	language, err := DetectByFile("../testdata/mock_file.js")

	if err != nil {
		t.Errorf("Should not fail to load test file")
	}

	if language != "JavaScript" {
		t.Errorf("Should detect JavasScript language")
	}

	language, err = DetectByFile("../testdata/invalid_file")

	if err != nil {
		t.Errorf("Should not fail to load invalid file")
	}

	if language != "" {
		t.Errorf("Should not detect any language")
	}
}
