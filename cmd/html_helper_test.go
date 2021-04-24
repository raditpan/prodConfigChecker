package cmd

import (
	"github.com/sergi/go-diff/diffmatchpatch"

	"testing"
)

func Test_DiffPrettyHtmlLeft(t *testing.T) {

	tests := []struct {
		name    string // The name of the test
		text1	string
		text2	string
		expected	string
		fixTab	bool
	}{
		{"Delete diff at the end", "test1", "test2", "<span>test</span><del style=\"background:#ffb5b5;\">1</del>", false},
		{"Delete diff at the start", "yestest", "test", "<del style=\"background:#ffb5b5;\">yes</del><span>test</span>", false},
		{"Delete diff in the middle", "test123text", "testtext", "<span>test</span><del style=\"background:#ffb5b5;\">123</del><span>text</span>", false},
		{"Fix space and tabs", " test1", " test2", "<span>&nbsp;test</span><del style=\"background:#ffb5b5;\">1</del>", true},
		{"Replace newline", "test\n1", "test\n2", "<span>test<br></span><del style=\"background:#ffb5b5;\">1</del>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(tt.text1, tt.text2, false)

			result := DiffPrettyHtmlLeft(diffs, tt.fixTab)

			if len(diffs) == 1 {
				t.Errorf("Result html diff has no diff")
			}

			if result != tt.expected {
				t.Errorf("Result html diff incorrect, got: %s", result)
			}
		})
	}
}

func Test_DiffPrettyHtmlRight(t *testing.T) {

	tests := []struct {
		name    string // The name of the test
		text1	string
		text2	string
		expected	string
		fixTab	bool
	}{
		{"Insert diff at the end", "test1", "test2", "<span>test</span><span style=\"background:#d1ffd1;\">2</span>", false},
		{"Insert diff at the start", "test", "yestest", "<span style=\"background:#d1ffd1;\">yes</span><span>test</span>", false},
		{"Insert diff in the middle", "testtext", "test123text", "<span>test</span><span style=\"background:#d1ffd1;\">123</span><span>text</span>", false},
		{"Fix space and tabs", " test1", " test2", "<span>&nbsp;test</span><span style=\"background:#d1ffd1;\">2</span>", true},
		{"Replace newline", "test\n1", "test\n2", "<span>test<br></span><span style=\"background:#d1ffd1;\">2</span>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(tt.text1, tt.text2, false)

			result := DiffPrettyHtmlRight(diffs, tt.fixTab)

			if len(diffs) == 1 {
				t.Errorf("Result html diff has no diff")
			}

			if result != tt.expected {
				t.Errorf("Result html diff incorrect, got: %s", result)
			}
		})
	}
}