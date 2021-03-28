# prodConfigChecker
Command line for checking diff between prod and qa config, to help reduce human error when working with config files in multiple environments.

# How to run
```
prodConfigChecker run <app_name>
```

# Build/run from source
First, install Go in your machine
```
brew install go
```

Go to the directory you clone the project to. You can run the app with these commands:
```
// get all the dependencies
go get -d ./...

// run the main app
go run main.go
```