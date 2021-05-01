package cmd

import (
	"io/ioutil"

	"strings"

	"bytes"

	"html"

	"strconv"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func writeHtmlFile(diffArray []ConfigDiffItem, appName string) string {

	var sb strings.Builder
	sb.WriteString("<h2>" + appName + " - config diff report &#128203;</h2>")
	sb.WriteString("<div> number of diff files : " + strconv.Itoa(len(diffArray)) + "</div>")
	sb.WriteString("<hr>");

		for _, htmlDiff := range diffArray {
			sb.WriteString("<div style=\"overflow: auto;\">")
			sb.WriteString("<h3>" + htmlDiff.fileName + " : </h3><br>")
			sb.WriteString("<div style=\"float: left;width: 48%; border-right: 2px solid #808080;\">")
			sb.WriteString("<b> QA</b><br><br>")
			sb.WriteString(htmlDiff.diffLeft)
			sb.WriteString("</div>")
			sb.WriteString("<div style=\"float: left;width: 50%; margin-left: 1em;\">");
			sb.WriteString("<b> PROD</b><br><br>")
			sb.WriteString(htmlDiff.diffRight)
			sb.WriteString("</div>")
			sb.WriteString("</div>")
			sb.WriteString("<hr>")
		}

	outputFileName := appName + "_config_diff.html"
	ioutil.WriteFile(outputFileName, []byte(sb.String()), 0644)
	
	return outputFileName
}

type ConfigDiffItem struct {
	fileName string
	diffLeft string
	diffRight string
}


func DiffPrettyHtmlLeft(diffs []diffmatchpatch.Diff, doFixTab bool) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)

		if doFixTab {
			text = strings.Replace(text, " ", "&nbsp;", -1)
		}

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("<del style=\"background:#ffb5b5;\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</del>")
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString("<span style=\"word-wrap:break-word\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		}
	}
	return buff.String()
}


func DiffPrettyHtmlRight(diffs []diffmatchpatch.Diff, doFixTab bool) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
		
		if doFixTab {
			text = strings.Replace(text, " ", "&nbsp;", -1)
		}

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("<span style=\"background:#d1ffd1;\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString("<span style=\"word-wrap:break-word\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		}
	}
	return buff.String()
}