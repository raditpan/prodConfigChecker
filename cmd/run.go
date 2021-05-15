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

	"io/fs"

	"fmt"

	"strings"

	"errors"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var repo string
var silentMode bool
var colorBlue = "\033[34m"
var colorReset = "\033[0m"

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Compare production config with qa config",
	Long: `Compare production config with qa config by specifying the name of application to compare
	 For example:

prodConfigChecker run <app name>
prodConfigChecker run <app name> --repo <absolute path to your config repo>`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires app name argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(string(colorBlue), "app name to check : "+args[0])

		appName := args[0]

		var configRepoPath string
		if len(repo) > 0 {
			// use cmd flag first, if available
			configRepoPath = repo
		} else {
			// else use from yaml file in home directory
			configRepoPath = viper.GetString("configRepoPath")
		}

		files := getFileListInDirectory(configRepoPath, "qa", appName)
		diffArray := diffConfigFiles(configRepoPath, appName, files, silentMode)

		outputFileName := writeHtmlFile(diffArray, appName)

		fmt.Println(string(colorBlue), "=====================================")
		fmt.Println("HTML output file : " + outputFileName)
	},
}

func getFileListInDirectory(configRepoPath string, envName string, appName string) []fs.FileInfo {
	files, err := ioutil.ReadDir(configRepoPath + "/" + envName + "/" + appName)
	filtered := []fs.FileInfo{}
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		// filter out system files with '.' as prefix
		if !strings.HasPrefix(f.Name(), ".") && !f.IsDir() {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

func getFileContent(configRepoPath string, envName string, appName string, fileName string) string {
	byteContent, err2 := ioutil.ReadFile(configRepoPath + "/" + envName + "/" + appName + "/" + fileName)

	if err2 != nil {
		return ""
	}

	return string(byteContent[:])
}

func diffConfigFiles(configRepoPath string, appName string, files []fs.FileInfo, silent bool) []ConfigDiffItem {

	diffArray := make([]ConfigDiffItem, 0)

	for _, f := range files {

		var item ConfigDiffItem

		qaFileString := getFileContent(configRepoPath, "qa", appName, f.Name())
		prodFileString := getFileContent(configRepoPath, "production", appName, f.Name())

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(qaFileString, prodFileString, false)

		if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
			// skip the file for when there's no diff case
			continue
		}

		fmt.Println(diffs)

		if !silent {
			fmt.Println(string(colorBlue), "=====================================")
			fmt.Println(string(colorBlue), f.Name()+" config files diff : ", string(colorReset))
			fmt.Println(dmp.DiffPrettyText(diffs))
		}

		shouldFixTab := isYamlFile(f.Name())

		item.fileName = f.Name()

		if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffDelete {
			item.diffLeft = SimpleDiffFormat(diffs[0])
			item.diffRight = "<span style=\"color:red\">No file available</span>"
		} else {
			item.diffLeft = DiffPrettyHtmlLeft(diffs, shouldFixTab)
			item.diffRight = DiffPrettyHtmlRight(diffs, shouldFixTab)
		}

		diffArray = append(diffArray, item)
	}

	return diffArray
}

func isYamlFile(fileName string) bool {
	fileExtension := filepath.Ext(fileName)
	return fileExtension == ".yml" || fileExtension == ".yaml"
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	rootCmd.PersistentFlags().StringVar(&repo, "repo", "", "Absolute path to your config repo")
	rootCmd.PersistentFlags().BoolVarP(&silentMode, "silent", "s", false, "Silence diff result in console output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
