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
var ecsRepoMode bool
var colorBlue = "\033[34m"
var colorReset = "\033[0m"
var qaFolder string
var prodFolder string

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

		if ecsRepoMode {
			qaFolder = "th/staging"
			prodFolder = "th/prod"
		} else {
			qaFolder = "qa"
			prodFolder = "production"
		}

		qaFiles := getFileListInDirectory(configRepoPath, qaFolder, appName)
		prodFiles := getFileListInDirectory(configRepoPath, prodFolder, appName)
		files := mergeFileList(qaFiles, prodFiles)
		diffArray := diffConfigFiles(configRepoPath, qaFolder, prodFolder, appName, files, silentMode)

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

func getFileContent(configRepoPath string, envName string, appName string, fileName string) (string, bool) {
	byteContent, err2 := ioutil.ReadFile(configRepoPath + "/" + envName + "/" + appName + "/" + fileName)

	if err2 != nil {
		fmt.Println(err2)
		return "", false
	}

	return string(byteContent[:]), true
}

func mergeFileList(first []fs.FileInfo, second []fs.FileInfo) []fs.FileInfo {
	for _, f := range first {
		exist := false
		for _, s := range second {
			if f.Name() == s.Name() {
				exist = true
				break
			}
		}
		if !exist {
			second = append(second, f)
		}
	}

	return second
}

func diffConfigFiles(configRepoPath string, qaFolder string, prodFolder string, appName string, files []fs.FileInfo, silent bool) []ConfigDiffItem {

	diffArray := make([]ConfigDiffItem, 0)

	for _, f := range files {

		var item ConfigDiffItem

		if !silent {
			fmt.Println(string(colorBlue), "=====================================")
		}

		qaFileString, qafileExist := getFileContent(configRepoPath, qaFolder, appName, f.Name())
		prodFileString, prodfileExist := getFileContent(configRepoPath, prodFolder, appName, f.Name())

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(qaFileString, prodFileString, false)

		if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
			// skip the file for when there's no diff case
			item.noDiff = true
		}

		if !silent {

			if item.noDiff {
				fmt.Println(string(colorBlue), f.Name()+" has no diff ", string(colorReset))
			} else {
				fmt.Println(string(colorBlue), f.Name()+" config files diff : ", string(colorReset))
				fmt.Println(dmp.DiffPrettyText(diffs))
			}
		}

		shouldFixTab := isYamlFile(f.Name())

		item.fileName = f.Name()
		noFileWarningSpan := "<span style=\"color:red\">⚠️ No file available</span>"

		if qafileExist && !prodfileExist && len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffDelete {
			// case where there's only file on QA, but not in Prod
			item.diffLeft = simpleDiffFormat(diffs[0])
			item.diffRight = noFileWarningSpan
		} else if !qafileExist && prodfileExist && len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffInsert {
			// case where there's only file on Prod, but not in QA
			item.diffLeft = noFileWarningSpan
			item.diffRight = simpleDiffFormat(diffs[0])
		} else if qafileExist && !prodfileExist && qaFileString == "" {
			// case where there's only file on QA, but not in Prod, and QA file is empty
			item.diffLeft = ""
			item.diffRight = noFileWarningSpan
		} else if !qafileExist && prodfileExist && prodFileString == "" {
			// case where there's only file on Prod, but not in QA, and Prod file is empty
			item.diffLeft = noFileWarningSpan
			item.diffRight = ""
		} else {
			item.diffLeft = diffPrettyHtmlLeft(diffs, shouldFixTab)
			item.diffRight = diffPrettyHtmlRight(diffs, shouldFixTab)
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
	rootCmd.PersistentFlags().BoolVarP(&ecsRepoMode, "ecs", "e", false, "Use ECS repo folder structure")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
