package cmd

import (
	"testing"
)

func Test_GetFileContentQa(t *testing.T) {

	result := getFileContent("../testdata", "qa", "acm-test", "application.yaml")

	if result != "testconfig: true" {
		t.Errorf("Result file content not as expected, got: %s", result)
	}
}

func Test_GetFileContentProd(t *testing.T) {

	result := getFileContent("../testdata", "production", "acm-test", "application.yaml")

	if result != "testconfig: false" {
		t.Errorf("Result file content not as expected, got: %s", result)
	}
}

func Test_GetFileListInDirectory(t *testing.T) {

	result := getFileListInDirectory("../testdata", "acm-test")

	if len(result) != 2 {
		t.Errorf("Incorrect number of files")
	}

	if result[0].Name() != "application.yaml" {
		t.Errorf("Incorrect file name, got: %s", result[0].Name())
	}

	if result[1].Name() != "config.json" {
		t.Errorf("Incorrect file name, got: %s", result[0].Name())
	}
}

func Test_IsYamlFile(t *testing.T) {
	tests := []struct {
		name    string // The name of the test
		fileName	string
		expected	bool
	}{
		{"Not yaml file", "application.json", false},
		{"yaml file", "application.yaml", true},
		{"yml file", "application.yml", true},
	}

	for _, tt := range tests {
		result := isYamlFile(tt.fileName)

		if result != tt.expected {
			t.Errorf("Result not as expected, got: %t", result)
		}
	}
}