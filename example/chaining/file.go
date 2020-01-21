package main

import "fmt"

// PrintAPIKey prints the API key stored in ENV
func PrintAPIKey() {
	fmt.Println(ENV.Get("API_KEY"))
}
