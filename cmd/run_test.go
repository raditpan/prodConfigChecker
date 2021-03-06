package cmd

import (
	"testing"
)

func Test_GetFileContentQa(t *testing.T) {

	result, fileExist := getFileContent("../testdata", "qa", "acm-test", "application.yaml")

	if result != "testconfig: true" || !fileExist {
		t.Errorf("Result file content not as expected, got: %s", result)
	}
}

func Test_GetFileContentProd(t *testing.T) {

	result, fileExist := getFileContent("../testdata", "production", "acm-test", "application.yaml")

	if result != "testconfig: false" || !fileExist {
		t.Errorf("Result file content not as expected, got: %s", result)
	}
}

func Test_GetFileContentQa_EmptyFile(t *testing.T) {

	result, fileExist := getFileContent("../testdata", "qa", "acm-test", "test.txt")

	if result != "" || !fileExist {
		t.Errorf("Result file content not as expected, got: %s", result)
	}
}

func Test_GetFileContent_NoFile_ReturnEmptyString(t *testing.T) {
	result, fileExist := getFileContent("../testdata", "production", "acm-test-2", "application.yaml")

	if result != "" || fileExist {
		t.Errorf("Result file content not as empty, got: %s", result)
	}
}

func Test_GetFileListInDirectory(t *testing.T) {

	result := getFileListInDirectory("../testdata", "production", "acm-test")

	if len(result) != 2 {
		t.Errorf("Incorrect number of files")
	}

	if result[0].fileInfo.Name() != "application.yaml" {
		t.Errorf("Incorrect file name, got: %s", result[0].fileInfo.Name())
	}

	if result[1].fileInfo.Name() != "config.json" {
		t.Errorf("Incorrect file name, got: %s", result[0].fileInfo.Name())
	}
}

func Test_GetFileListInDirectory_IncludeInnerDirectory(t *testing.T) {

	result := getFileListInDirectory("../testdata", "qa", "acm-test")

	if len(result) != 5 {
		t.Errorf("Incorrect number of files")
	}

	exist := false
	for _, f := range result {
		if f.fileInfo.Name() == "rules.js" {
			exist = true

			if f.parentFolder != "validation-rules" {
				t.Errorf("parent folder name not correct")
			}
		}
	}

	if !exist {
		t.Errorf("sub-folder file not included")
	}
}

func Test_GetFileListInDirectory_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	getFileListInDirectory("../testdata", "production", "acm-test-2")
}

func Test_MergeFileList(t *testing.T) {
	qaFiles := getFileListInDirectory("../testdata", "qa", "acm-test")
	prodFiles := getFileListInDirectory("../testdata", "production", "acm-test")

	result := mergeFileList(qaFiles, prodFiles)

	if len(result) != 5 {
		t.Errorf("Incorrect number of files")
	}

	exist := false
	for _, f := range result {
		if f.fileInfo.Name() == "qa-application.yaml" {
			exist = true
		}
	}

	if !exist {
		t.Errorf("qa-only config not exist after merging")
	}
}

func Test_DiffConfigFiles(t *testing.T) {

	files := getFileListInDirectory("../testdata", "production", "acm-test")
	result := diffConfigFiles("../testdata", "qa", "production", "acm-test", files, true)

	if len(result) != 2 {
		t.Errorf("Incorrect number of files")
	}

	if result[0].fileName != "application.yaml" {
		t.Errorf("Incorrect file name, got: %s", result[0].fileName)
	}

	if result[0].diffLeft != "<span style=\"word-wrap:break-word\">testconfig:&nbsp;</span><del style=\"background:#ffb5b5;word-wrap:break-word\">tru</del><span style=\"word-wrap:break-word\">e</span>" {
		t.Errorf("Incorrect diff left, got: %s", result[0].diffLeft)
	}

	if result[0].diffRight != "<span style=\"word-wrap:break-word\">testconfig:&nbsp;</span><span style=\"background:#d1ffd1;word-wrap:break-word\">fals</span><span style=\"word-wrap:break-word\">e</span>" {
		t.Errorf("Incorrect diff left, got: %s", result[0].diffRight)
	}

	if result[1].fileName != "config.json" {
		t.Errorf("Incorrect file name, got: %s", result[1].fileName)
	}

	if !result[1].noDiff {
		t.Errorf("Unexpected diff in config.json")
	}
}

func Test_DiffConfigFiles_NonSilent(t *testing.T) {

	files := getFileListInDirectory("../testdata", "production", "acm-test")
	result := diffConfigFiles("../testdata", "qa", "production", "acm-test", files, false)

	if len(result) != 2 {
		t.Errorf("Incorrect number of files")
	}
}

func Test_IsYamlFile(t *testing.T) {
	tests := []struct {
		name     string // The name of the test
		fileName string
		expected bool
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
