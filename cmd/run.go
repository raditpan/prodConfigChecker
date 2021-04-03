/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"io/ioutil"

	"fmt"

	"strings"

	"bytes"

	"html"

	"strconv"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Compare production config with qa config",
	Long: `Compare production config with qa config by specifying the name of application to compare
	 For example:

prodConfigChecker run acm-bpay-api`,
	Run: func(cmd *cobra.Command, args []string) {
		colorBlue := "\033[34m"
		colorReset := "\033[0m"

		fmt.Println(string(colorBlue), "app name to check : " + args[0])

		appName := args[0]
		configRepoPath := viper.GetString("configRepoPath")


		files, err := ioutil.ReadDir(configRepoPath + "/production/" + appName)
		if err != nil {
			panic(err)
		}

		diffArray := make([]ConfigDiffItem, 0)

		for _, f := range files {
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}

			var item ConfigDiffItem
			item.fileName = f.Name()

			prod, err := ioutil.ReadFile(configRepoPath + "/production/" + appName + "/" + f.Name())
			if err != nil{
				panic(err)
			}
			qa, err2 := ioutil.ReadFile(configRepoPath + "/qa/" + appName + "/" + f.Name())
			if err2 != nil{
				panic(err2)
			}
			qaFileString := string(qa[:])
			prodFileString := string(prod[:])
			dmp := diffmatchpatch.New()

			diffs := dmp.DiffMain(qaFileString, prodFileString, false)

			if len(diffs) == 1 {
				// skip the file for no diff case
				continue
			}

			dmp.DiffCleanupSemantic(diffs)

			fmt.Println(string(colorBlue), "=====================================")
			fmt.Println(string(colorBlue), f.Name() + " config files diff : ", string(colorReset))
			fmt.Println(dmp.DiffPrettyText(diffs))

			// var patchList = dmp.PatchMake(qaFileString, prodFileString, diffs)
			// var patchText = dmp.PatchToText(patchList)
			// fmt.Println(patchText)
			fileExtension := filepath.Ext(f.Name())
			shouldFixTab := fileExtension == ".yml" || fileExtension == ".yaml"

			item.diffLeft = DiffPrettyHtmlLeft(diffs, shouldFixTab)
			item.diffRight = DiffPrettyHtmlRight(diffs, shouldFixTab)

			diffArray = append(diffArray, item)
		}

		writeHtmlFile(diffArray, appName)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func writeHtmlFile(diffArray []ConfigDiffItem, appName string) {

	var sb strings.Builder
	sb.WriteString("<h2>" + appName + " - config diff report &#128203;</h2>")
	sb.WriteString("<div> number of diff files : " + strconv.Itoa(len(diffArray)) + "</div>")
	sb.WriteString("<hr>");

		for _, htmlDiff := range diffArray {
			sb.WriteString("<div style=\"overflow: auto;\">")
			sb.WriteString("<h3>" + htmlDiff.fileName + " : </h3><br><br>")
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

	ioutil.WriteFile(appName + "_config_diff.html", []byte(sb.String()), 0644)
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
			_, _ = buff.WriteString("<span>")
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
			_, _ = buff.WriteString("<ins style=\"background:#d1ffd1;\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</ins>")
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString("<span>")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		}
	}
	return buff.String()
}