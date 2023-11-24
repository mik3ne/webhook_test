package main

import (
	"webhook/internal/app"
	"webhook/internal/config"
)

func main() {

	//
	// using config injection for better testing facilities
	//
	app.CreateApp(config.LoadConfig).Run()
}
