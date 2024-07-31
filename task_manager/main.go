package main

import "task_manager/router"

func main() {
	router := router.InitializeRouter()
	router.Run(":8080")
}
