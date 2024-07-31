package main

import "task_manager/router"

// main is the entry point of the program.
// It initializes the router and starts the server on port 8080.
func main() {
	router := router.InitializeRouter()
	router.Run(":8080")
}
