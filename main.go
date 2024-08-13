package main

import (
	"lib_mngmt/controllers"
	"lib_mngmt/services"
)

func main() {
	library := services.NewLibrary()
	controller := controllers.NewLibraryController(library)
	controller.Run()
}
