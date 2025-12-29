package main

import (
	"log"
	"os"

	"github.com/manosriram/wingman/internal/repository"
)

func main() {
	// TODO: use wd
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting WorkingDir")
	}
	targetDir := wd

	err = repository.NewRepository(targetDir).Run()
	if err != nil {
		log.Fatalf("Error initializing program: %s\n", err.Error())
	}
}
