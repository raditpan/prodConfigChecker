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