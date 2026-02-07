package main

import (
	"ignis/app"
	"log"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Printf("failed to run app: %v", err)
	}

	log.Println("Application exited")
}
