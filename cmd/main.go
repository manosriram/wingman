package main

import (
	"log"

	"github.com/manosriram/wingman/internal/repository"
)

func main() {
	targetDir := "/Users/manosriram/go/src/go2java"

	err := repository.NewRepository(targetDir).Run()
	if err != nil {
		log.Fatalf("Error initializing program\n")
	}
}
