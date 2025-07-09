package main

import (
	"media-service/pkg/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}