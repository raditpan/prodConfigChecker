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