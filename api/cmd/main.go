package main

import (
	// "media-service/cmd"
	authApp "auth-service/pkg/app"
	// "notes-service/cmd"
)

func main() {
	if err := authApp.Run(); err != nil {
		panic(err)
	}
}