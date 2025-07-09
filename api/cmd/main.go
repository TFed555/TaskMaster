package main

import (
	mediaApp "media-service/pkg/app"
	authApp "auth-service/pkg/app"
	notesApp "notes-service/pkg/app"
)

func main() {
	go func() {
		if err := mediaApp.Run(); err != nil {
			panic(err)
		}
	}()
	go func() {
		if err := notesApp.Run(); err != nil {
			panic(err)
		}
	}()

	if err := authApp.Run(); err != nil {
		panic(err)
	}

	
}