package main

import (
	"fmt"
	golangdotenv "github.com/navaz-alani/golang-dotenv"
)

// ENV stores the program's environment variables
var ENV golangdotenv.Env
var envImportErr error

func main() {
	// Initialize the environment...
	ENV, envImportErr = golangdotenv.Load(".env", true)

	if envImportErr != nil {
		fmt.Print(envImportErr)
	}

	// Print API key from other file
	PrintAPIKey()
}
