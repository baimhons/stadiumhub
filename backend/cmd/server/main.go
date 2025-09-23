package main

import "github.com/baimhons/stadiumhub/internal/initial"

func main() {
	app := initial.InitializeApp()

	app.Run()

}
