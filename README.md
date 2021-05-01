# prodConfigChecker

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-63%25-brightgreen.svg?longCache=true&style=flat)</a>

Command line tool for checking diff between prod and qa config, to help reduce human error when working with config files in multiple environments.

Config repo's folder structure that is supported:
```
config-repo
├── production
│   └── <app-name>
│       └── <config files>
├── qa
│   └── <app-name>
│       └── <config files>
```

## Prepare Go enviroment
These are steps to setup your environment to install and run Go package easily.
First, install Go in your machine:
```
brew install go
```

Add below GOPATH var in your ~/.zshrc or ~/.bash_profile:
```
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=${PATH}:$GOBIN
```

## How to run
Install the command line tool:
```
go get github.com/raditpan/prodConfigChecker
```


Create `.prodConfigChecker.yaml` file in your home directoy with path to your local config repo:
```
configRepoPath: "<absolute path to your config repo>"
```

Run the command:

```
prodConfigChecker run <app_name>

// run with repo path option
prodConfigChecker run <app_name> --repo <absolute path to your config repo>

// run with custom config file
prodConfigChecker run <app_name> --config custom-config.yaml

// run with silent mode, no files diff in console output
prodConfigChecker run <app_name> -s
```

Check the diff output in terminal console. HTML output file is also generated in your current directory, in case you want to share the result with others.

## Build/run from source

Go to the directory you clone the project to. You can run the app with these commands:
```
// get all the dependencies
go get -d ./...

// run the main app
go run main.go run <app_name>
```