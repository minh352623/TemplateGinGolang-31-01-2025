package main

import "ecom/internal/inittiallize"

// @securityDefinitions.apiKey bearerToken
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
func main() {
	inittiallize.Run()
}
