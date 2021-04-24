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

	"errors"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var repo string
var silentMode bool

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
		colorBlue := "\033[34m"
		colorReset := "\033[0m"

		fmt.Println(string(colorBlue), "app name to check : " + args[0])

		appName := args[0]

		var configRepoPath string
		if len(repo) > 0 {
			configRepoPath = repo
		} else {
			configRepoPath = viper.GetString("configRepoPath")
		}

		files, err := ioutil.ReadDir(configRepoPath + "/production/" + appName)
		if err != nil {
			panic(err)
		}

		diffArray := make([]ConfigDiffItem, 0)

		for _, f := range files {
			if strings.HasPrefix(f.Name(), ".") {
				// skip system files with '.' as prefix
				continue
			}

			var item ConfigDiffItem

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
				// skip the file for when there's no diff case
				continue
			}

			if(!silentMode) {
				fmt.Println(string(colorBlue), "=====================================")
				fmt.Println(string(colorBlue), f.Name() + " config files diff : ", string(colorReset))
				fmt.Println(dmp.DiffPrettyText(diffs))
			}

			fileExtension := filepath.Ext(f.Name())
			shouldFixTab := fileExtension == ".yml" || fileExtension == ".yaml"

			item.fileName = f.Name()
			item.diffLeft = DiffPrettyHtmlLeft(diffs, shouldFixTab)
			item.diffRight = DiffPrettyHtmlRight(diffs, shouldFixTab)

			diffArray = append(diffArray, item)
		}

		outputFileName := writeHtmlFile(diffArray, appName)

		fmt.Println(string(colorBlue), "=====================================")
		fmt.Println("HTML output file : " + outputFileName)
	},
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
