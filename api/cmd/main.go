package main

import (
	// mediaApp "media-service/pkg/app"
	authApp "auth-service/pkg/app"
	// "notes-service/cmd"
)

func main() {
	// if err := mediaApp.Run(); err != nil {
	// 	panic(err)
	// }

	if err := authApp.Run(); err != nil {
		panic(err)
	}
}