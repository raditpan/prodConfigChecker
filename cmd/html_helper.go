package cmd

import (
	"io/ioutil"

	"strings"

	"bytes"

	"html"

	"strconv"

	"github.com/sergi/go-diff/diffmatchpatch"

	"time"
)

func writeHtmlFile(diffArray []ConfigDiffItem, appName string, qaTitle string) string {

	var sb strings.Builder
	currentTime := time.Now()

	var diffCount int
	for _, htmlDiff := range diffArray {
		if !htmlDiff.noDiff {
			diffCount++
		}
	}

	sb.WriteString("<h2>" + appName + " - config diff report &#128203;</h2>")
	sb.WriteString("<div> Run date-time : " + currentTime.Format("02-Jan-2006 15:04:05") + "</div>")
	sb.WriteString("<div> Number of files : " + strconv.Itoa(len(diffArray)) + "</div>")
	sb.WriteString("<div> Number of diff files : " + strconv.Itoa(diffCount) + "</div>")
	sb.WriteString("<hr>")

	for _, htmlDiff := range diffArray {
		if htmlDiff.noDiff {
			sb.WriteString("<div style=\"overflow: auto;\">")
			sb.WriteString("<h3> ✅ " + htmlDiff.fileName + " - no diff</h3>")
			sb.WriteString("</div>")
			sb.WriteString("<hr>")
			continue
		}

		sb.WriteString("<div style=\"overflow: auto;\">")
		sb.WriteString("<h3> 🛂 " + htmlDiff.fileName + " : </h3><br>")
		sb.WriteString("<div style=\"float: left;width: 48%; border-right: 2px solid #808080;\">")
		sb.WriteString("<b> " + qaTitle + "</b><br><br>")
		sb.WriteString(htmlDiff.diffLeft)
		sb.WriteString("</div>")
		sb.WriteString("<div style=\"float: left;width: 50%; margin-left: 1em;\">")
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
	fileName  string
	diffLeft  string
	diffRight string
	noDiff    bool
}

func simpleDiffFormat(diff diffmatchpatch.Diff) string {
	return strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
}

func diffPrettyHtmlLeft(diffs []diffmatchpatch.Diff, doFixTab bool) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)

		if doFixTab || strings.TrimSpace(text) == "" {
			text = strings.Replace(text, " ", "&nbsp;", -1)
		}

		if text == "<br>" {
			// add space to hilight line removal, because lone <br> inside span will not get background
			text = "<br>&nbsp;"
		}

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("<del style=\"background:#ffb5b5;word-wrap:break-word\">")
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

func diffPrettyHtmlRight(diffs []diffmatchpatch.Diff, doFixTab bool) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)

		if doFixTab || strings.TrimSpace(text) == "" {
			text = strings.Replace(text, " ", "&nbsp;", -1)
		}

		if text == "<br>" {
			// add space to hilight line addition, because lone <br> inside span will not get background
			text = "<br>&nbsp;"
		}

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("<span style=\"background:#d1ffd1;word-wrap:break-word\">")
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
