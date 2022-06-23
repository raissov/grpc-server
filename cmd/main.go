package main

import (
	"os"
	"server/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
